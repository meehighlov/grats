package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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

	HEADER_MESSAGE_LIST_NOT_EMPTY = "Добавленные дни рождения✨"
	HEADER_MESSAGE_LIST_NOT_EMPTY_CHAT = "Список др в чате %s✨"
	HEADER_MESSAGE_LIST_IS_EMPTY  = "Записей пока нет✨"
)

func ListBirthdaysHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	message := event.GetMessage()

	chatId := message.Chat.Id
	if event.GetCallbackQuery().Id != "" {
		chatIdStr := common.CallbackFromString(event.GetCallbackQuery().Data).Id
		if chatId_, err_ := strconv.Atoi(chatIdStr); err_ != nil {
			event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
			return err_
		} else {
			chatId = chatId_
		}	
	}

	friends, err := (&db.Friend{UserId: message.From.Id, ChatId: chatId}).Filter(ctx, tx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return err
	}

	if event.GetCallbackQuery().Id != "" {
		event.EditCalbackMessage(
			ctx,
			buildChatHeaderMessage(ctx, chatId, event, (len(friends) == 0)),
			buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, strconv.Itoa(chatId)),
		)

		return nil
	}

	event.ReplyWithKeyboard(
		ctx,
		buildChatHeaderMessage(ctx, chatId, event, (len(friends) == 0)),
		buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET, strconv.Itoa(chatId)),
	)

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

func buildPagiButtons(total, limit, offset int, chatId string) [][]map[string]string {
	if total == 0 {
		return [][]map[string]string{}
	}
	if offset == total {
		return [][]map[string]string{{
			{
				"text":          "свернуть👆",
				"callback_data": common.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<", chatId).String(),
			},
		}}
	}
	var keyBoard []map[string]string
	if offset+limit >= total {
		previousButton := map[string]string{"text": "👈назад", "callback_data": common.CallList(strconv.Itoa(offset), "<<", chatId).String()}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": "вперед👉", "callback_data": common.CallList(strconv.Itoa(offset), ">>", chatId).String()}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": "вперед👉", "callback_data": common.CallList(strconv.Itoa(offset), ">>", chatId).String()}
			previousButton := map[string]string{"text": "👈назад", "callback_data": common.CallList(strconv.Itoa(offset), "<<", chatId).String()}
			keyBoard = []map[string]string{previousButton, nextButton}
		}
	}

	allButton := map[string]string{"text": fmt.Sprintf("показать все (%d)👇", total), "callback_data": common.CallList(strconv.Itoa(offset), "<>", chatId).String()}
	allButtonBar := []map[string]string{allButton}

	markup := [][]map[string]string{}
	if total <= limit {
		return markup
	}

	markup = append(markup, keyBoard)
	markup = append(markup, allButtonBar)

	return markup
}

func ListPaginationCallbackQueryHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	chatId, err := strconv.Atoi(params.BoundChat)
	if err != nil {
		return err
	}

	friends, err := (&db.Friend{UserId: callbackQuery.From.Id, ChatId: chatId}).Filter(ctx, tx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return err
	}

	direction := params.Pagination.Direction

	slog.Debug(fmt.Sprintf("direction: %s limit: %d offset: %s", direction, limit_, offset))

	if direction == "<" {
		slog.Debug("back to previous screen, offset not changed")
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

	msg := buildChatHeaderMessage(ctx, chatId, event, (len(friends) == 0))

	event.EditCalbackMessage(ctx, msg, buildFriendsListMarkup(friends, limit_, offset_, params.BoundChat))

	return nil
}

func buildFriendsButtons(friends []db.Friend, limit, offset int, callbackDataBuilder func(id string, offset int) string) []map[string]string {
	sort.Slice(friends, func(i, j int) bool { return birthdayComparator(friends, i, j) })
	var buttons []map[string]string
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

		button := map[string]string{
			"text":          buttonText,
			"callback_data": callbackDataBuilder(friend.ID, offset),
		}
		buttons = append(buttons, button)
	}

	return buttons
}

func buildFriendsListMarkup(friends []db.Friend, limit, offset int, chatId string) [][]map[string]string {
	callbackDataBuilder := func(id string, offset int) string {
		return common.CallInfo(id, strconv.Itoa(offset), chatId).String()
	}
	friendsListAsButtons := buildFriendsButtons(friends, limit, offset, callbackDataBuilder)
	pagiButtons := buildPagiButtons(len(friends), limit, offset, chatId)

	markup := [][]map[string]string{}

	for _, button := range friendsListAsButtons {
		markup = append(markup, []map[string]string{button})
	}

	markup = append(markup, pagiButtons...)

	if strings.Contains(chatId, "-") {
		backToGroupButton := map[string]string{
			"text":          "👈 к чату",
			"callback_data": common.CallChatInfo(chatId).String(),
		}
		markup = append(markup, []map[string]string{backToGroupButton})
	}

	return markup
}

func buildChatHeaderMessage(ctx context.Context, chatId int, event common.Event, emptyList bool) string {
	if emptyList {
		return HEADER_MESSAGE_LIST_IS_EMPTY
	}
	chatFullInfo := event.GetChat(ctx, strconv.Itoa(chatId))
	if chatFullInfo.Id < 0 {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY_CHAT, chatFullInfo.Title)
	} else {
		return HEADER_MESSAGE_LIST_NOT_EMPTY
	}
}
