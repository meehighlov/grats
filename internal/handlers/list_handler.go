package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/models"
	"github.com/meehighlov/grats/telegram"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT = 5
	LIST_START_OFFSET = 0

	HEADER_MESSAGE_LIST_NOT_EMPTY = "Нажми, чтобы узнать детали✨"
	HEADER_MESSAGE_LIST_IS_EMPTY = "Записей пока нет✨"
)

func ListBirthdaysHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()
	friends, err := (&db.Friend{UserId: message.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return nil
	}

	if len(friends) == 0 {
		event.Reply(ctx, HEADER_MESSAGE_LIST_IS_EMPTY)
		return nil
	}

	event.ReplyWithKeyboard(
		ctx,
		HEADER_MESSAGE_LIST_NOT_EMPTY,
		buildFriendsListMarkup(friends, LIST_LIMIT, LIST_START_OFFSET),
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

func buildPagiButtons(total, limit, offset int) [][]map[string]string {
	if total == 0 {
		return [][]map[string]string{}
	}
	if offset == total {
		return [][]map[string]string{{
			{
				"text": "свернуть👆",
				"callback_data": models.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<").String(),
			},
		}}
	}
	var keyBoard []map[string]string
	if offset + limit >= total {
		previousButton := map[string]string{"text": "👈назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<").String()}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": "вперед👉", "callback_data": models.CallList(strconv.Itoa(offset), ">>").String()}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": "вперед👉", "callback_data": models.CallList(strconv.Itoa(offset), ">>").String()}
			previousButton := map[string]string{"text": "👈назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<").String()}
			keyBoard = []map[string]string{previousButton, nextButton}
		}
	}

	allButton := map[string]string{"text": fmt.Sprintf("показать все (%d)👇", total), "callback_data": models.CallList(strconv.Itoa(offset), "<>").String()}
	allButtonBar := []map[string]string{allButton}

	markup := [][]map[string]string{}
	if total <= limit {
		return markup
	}

	markup = append(markup, keyBoard)
	markup = append(markup, allButtonBar)

	return markup
}

func ListBirthdaysCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()
	callbackQuery := event.GetCallbackQuery()

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	friends, err := (&db.Friend{UserId: callbackQuery.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return nil
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

	msg := HEADER_MESSAGE_LIST_NOT_EMPTY
	if len(friends) == 0 {
		msg = HEADER_MESSAGE_LIST_IS_EMPTY
	}

	event.EditCalbackMessage(ctx, msg, buildFriendsListMarkup(friends, limit_, offset_))

	return nil
}

func buildFriendsButtons(friends []db.Friend, limit, offset int) []map[string]string {
	sort.Slice(friends, func(i, j int) bool { return birthdayComparator(friends, i, j) })
	var buttons []map[string]string
	for i, friend := range friends {
		if offset != len(friends) {
			if i == limit + offset {
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
			"text": buttonText,
			"callback_data": models.CallInfo(friend.ID, strconv.Itoa(offset)).String(),
		}
		buttons = append(buttons, button)
	}

	return buttons
}

func buildFriendsListMarkup(friends []db.Friend, limit, offset int) [][]map[string]string {
	friendsListAsButtons := buildFriendsButtons(friends, limit, offset)
	pagiButtons := buildPagiButtons(len(friends), limit, offset)

	markup := [][]map[string]string{}

	for _, button := range friendsListAsButtons {
		markup = append(markup, []map[string]string{button})
	}

	markup = append(markup, pagiButtons...)

	return markup
}
