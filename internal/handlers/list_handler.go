package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT            = 5
	LIST_START_OFFSET     = 0

	HEADER_MESSAGE_LIST_NOT_EMPTY      = "✨Личные напоминания о др"
	HEADER_MESSAGE_LIST_NOT_EMPTY_CHAT = "✨Список др в чате %s"
	HEADER_MESSAGE_LIST_IS_EMPTY       = "✨Записей пока нет"
)

func ListBirthdaysHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()

	tgChatId := message.GetChatIdStr()
	if event.GetCallbackQuery().Id != "" {
		tgChatId = common.CallbackFromString(event.GetCallbackQuery().Data).Id
	}

	// Сначала получаем чат по TGChatId
	chat := db.Chat{
		TGChatId: tgChatId,
	}
	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("Error fetching chats: " + err.Error())
		return err
	}

	var friends []db.Friend

	if len(chats) > 0 {
		// Получаем друзей по ID чата
		friends, err = (&db.Friend{ChatId: chats[0].ID}).Filter(ctx, tx)
		if err != nil {
			event.Logger.Error("Error fetching friends: " + err.Error())
			return err
		}
	} else {
		friends = []db.Friend{}
	}

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			buildChatHeaderMessage(ctx, tgChatId, event, (len(friends) == 0)),
			*buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, tgChatId).Murkup(),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		buildChatHeaderMessage(ctx, tgChatId, event, (len(friends) == 0)),
		*buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, tgChatId).Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func birthdayComparator(friends []db.Friend, i, j int) bool {
	if friends[i].IsTodayBirthday() {
		return true
	}
	if friends[j].IsTodayBirthday() {
		return false
	}
	countI := friends[i].CountDaysToBirthday()
	countJ := friends[j].CountDaysToBirthday()
	return countI > countJ
}

func appendControlButtons(keyboard *common.InlineKeyboard, total, limit, offset int, chatId string) error {
	buttons := []common.Button{}
	if total <= limit || total == 0 {
		return nil
	}
	if offset == total {
		buttons = append(buttons, *common.NewButton("⬆️", common.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<", chatId).String()))
		keyboard.AppendAsLine(buttons...)
		return nil
	}
	if offset+limit >= total {
		buttons = append(buttons, *common.NewButton("⬅️", common.CallList(strconv.Itoa(offset), "<<", chatId).String()))
	} else {
		if offset == 0 {
			buttons = append(buttons, *common.NewButton("➡️", common.CallList(strconv.Itoa(offset), ">>", chatId).String()))
		} else {
			buttons = append(buttons, *common.NewButton("⬅️", common.CallList(strconv.Itoa(offset), "<<", chatId).String()))
			buttons = append(buttons, *common.NewButton("➡️", common.CallList(strconv.Itoa(offset), ">>", chatId).String()))
		}
	}

	keyboard.AppendAsLine(buttons...)
	keyboard.AppendAsStack(*common.NewButton(fmt.Sprintf("(%d)⬇️", total), common.CallList(strconv.Itoa(offset), "<>", chatId).String()))

	return nil
}

func ListPaginationCallbackQueryHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		event.Logger.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	tgChatId := params.BoundChat

	// Сначала получаем чат по TGChatId
	chat := db.Chat{
		TGChatId: tgChatId,
	}
	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("Error fetching chats: " + err.Error())
		return err
	}

	var friends []db.Friend

	if len(chats) > 0 {
		// Получаем друзей по ID чата
		friends, err = (&db.Friend{ChatId: chats[0].ID}).Filter(ctx, tx)
		if err != nil {
			event.Logger.Error("Error fetching friends: " + err.Error())
			return err
		}
	} else {
		friends = []db.Friend{}
	}

	direction := params.Pagination.Direction

	event.Logger.Debug(fmt.Sprintf("direction: %s limit: %d offset: %s", direction, limit_, offset))

	if direction == "<" {
		event.Logger.Debug("back to previous screen, offset not changed")
	}
	if direction == "<<<" {
		offset_ = 0
	}
	if direction == ">>" {
		offset_ += LIST_PAGINATION_SHIFT
	}
	if direction == "<<" {
		offset_ -= LIST_PAGINATION_SHIFT
	}
	if direction == "<>" {
		offset_ = len(friends)
	}

	msg := buildChatHeaderMessage(ctx, tgChatId, event, (len(friends) == 0))

	if _, err := event.EditCalbackMessage(ctx, msg, *buildFriendsListMarkup(friends, limit_, offset_, tgChatId).Murkup()); err != nil {
		return err
	}

	return nil
}

func buildFriendsButtons(friends []db.Friend, limit, offset int, callbackDataBuilder func(id string, offset int) string) *[]common.Button {
	sort.Slice(friends, func(i, j int) bool { return birthdayComparator(friends, i, j) })
	buttons := []common.Button{}
	for i, friend := range friends {
		if offset != len(friends) {
			if i == limit+offset {
				break
			}
			if i < offset {
				continue
			}
		}

		buttonText := fmt.Sprintf("%s %s", friend.Name, friend.BirthDay)

		if friend.IsTodayBirthday() {
			buttonText = fmt.Sprintf("%s 🥳", buttonText)
		} else {
			if friend.IsThisMonthAfterToday() {
				buttonText = fmt.Sprintf("%s 🕒", buttonText)
			}
		}

		buttons = append(buttons, *common.NewButton(buttonText, callbackDataBuilder(friend.ID, offset)))
	}

	return &buttons
}

func buildFriendsListMarkup(friends []db.Friend, limit, offset int, chatId string) *common.InlineKeyboard {
	callbackDataBuilder := func(id string, offset int) string {
		return common.CallInfo(id, strconv.Itoa(offset)).String()
	}
	friendsListAsButtons := buildFriendsButtons(friends, limit, offset, callbackDataBuilder)
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsLine(*common.NewButton("🏠 в начало", common.CallSetup().String()))
	keyboard.AppendAsLine(*common.NewButton("➕ Добавить напоминание", common.CallAdd(chatId).String()))

	keyboard.AppendAsStack(*friendsListAsButtons...)

	appendControlButtons(keyboard, len(friends), limit, offset, chatId)

	if strings.Contains(chatId, "-") {
		keyboard.AppendAsLine(*common.NewButton("⬅️к чату", common.CallChatInfo(chatId).String()))
	}

	return keyboard
}

func buildChatHeaderMessage(ctx context.Context, chatId string, event *common.Event, emptyList bool) string {
	if emptyList {
		return HEADER_MESSAGE_LIST_IS_EMPTY
	}
	chatFullInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		return HEADER_MESSAGE_LIST_NOT_EMPTY_CHAT
	}
	if chatFullInfo.Id < 0 {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY_CHAT, chatFullInfo.Title)
	} else {
		return HEADER_MESSAGE_LIST_NOT_EMPTY
	}
}
