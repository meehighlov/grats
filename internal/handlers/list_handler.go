package handlers

import (
	"context"
	"database/sql"
	"fmt"
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

	chatId := message.GetChatIdStr()
	userId := strconv.Itoa(message.From.Id)
	if event.GetCallbackQuery().Id != "" {
		chatId = common.CallbackFromString(event.GetCallbackQuery().Data).Id
		userId = strconv.Itoa(event.GetCallbackQuery().From.Id)
	}

	totalCount, err := (&db.Friend{UserId: userId, ChatId: chatId}).Count(ctx, tx)
	if err != nil {
		event.Logger.Error("Error counting friends: " + err.Error())
		return err
	}

	friends, err := (&db.Friend{UserId: userId, ChatId: chatId}).Filter(ctx, tx,
		db.WithTodayBirthdaysFirst(),
		db.WithPagination(LIST_LIMIT, LIST_START_OFFSET))

	if err != nil {
		event.Logger.Error("Error fetching friends" + err.Error())
		return err
	}

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			buildChatHeaderMessage(ctx, chatId, event, totalCount == 0),
			*buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, chatId, totalCount).Murkup(),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		buildChatHeaderMessage(ctx, chatId, event, totalCount == 0),
		*buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, chatId, totalCount).Murkup(),
	); err != nil {
		return err
	}

	return nil
}

// birthdayComparator больше не используется, так как сортировка выполняется на стороне БД
// через функцию WithTodayBirthdaysFirst. Оставлена для справки.
func birthdayComparator(friends []*db.Friend, i, j int) bool {
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
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		event.Logger.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	chatId := params.BoundChat

	totalCount, err := (&db.Friend{UserId: strconv.Itoa(callbackQuery.From.Id), ChatId: chatId}).Count(ctx, tx)
	if err != nil {
		event.Logger.Error("Error counting friends: " + err.Error())
		return err
	}

	friends, err := (&db.Friend{UserId: strconv.Itoa(callbackQuery.From.Id), ChatId: chatId}).Filter(ctx, tx,
		db.WithTodayBirthdaysFirst(),
		db.WithPagination(limit_, offset_))

	if err != nil {
		event.Logger.Error("Error fetching friends" + err.Error())
		return err
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
		offset_ = totalCount
	}

	msg := buildChatHeaderMessage(ctx, chatId, event, (len(friends) == 0))

	if _, err := event.EditCalbackMessage(ctx, msg, *buildFriendsListMarkup(friends, limit_, offset_, chatId, totalCount).Murkup()); err != nil {
		return err
	}

	return nil
}

func buildFriendsButtons(friends []*db.Friend, _, offset int, callbackDataBuilder func(id string, offset int) string) *[]common.Button {
	buttons := []common.Button{}
	for _, friend := range friends {
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

func buildFriendsListMarkup(friends []*db.Friend, limit, offset int, chatId string, totalCount int) *common.InlineKeyboard {
	callbackDataBuilder := func(id string, offset int) string {
		return common.CallInfo(id, strconv.Itoa(offset)).String()
	}
	friendsListAsButtons := buildFriendsButtons(friends, limit, offset, callbackDataBuilder)
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsLine(*common.NewButton("🏠 в начало", common.CallSetup().String()))
	keyboard.AppendAsLine(*common.NewButton("➕ добавить напоминание", common.CallAddToChat(chatId).String()))

	keyboard.AppendAsStack(*friendsListAsButtons...)

	appendControlButtons(keyboard, totalCount, limit, offset, chatId)

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
