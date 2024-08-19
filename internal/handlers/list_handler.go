package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
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
	if offset == total {
		return [][]map[string]string{}
	}
	var keyBoard []map[string]string
	if offset + limit >= total {
		previousButton := map[string]string{"text": "назад", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:<<", limit, offset)}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": "вперед", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:>>", limit, offset)}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": "вперед", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:>>", limit, offset)}
			previousButton := map[string]string{"text": "назад", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:<<", limit, offset)}
			keyBoard = []map[string]string{previousButton, nextButton}
		}
	}

	allButton := map[string]string{"text": fmt.Sprintf("показать все (%d)", total), "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:<>", limit, offset)}
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
	payload := callbackQuery.Data
	params := strings.Split(payload, ";")
	limit := strings.Split(params[1], ":")[1]
	offset := strings.Split(params[2], ":")[1]

	limit_, err  := strconv.Atoi(limit)
	if err != nil {
		slog.Error("error parsing params in list pagination callback query: " + err.Error())
		return err
	}
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing params in list pagination callback query: " + err.Error())
		return err
	}

	friends, err := (&db.Friend{UserId: callbackQuery.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return nil
	}

	direction := strings.Split(strings.Split(callbackQuery.Data, ";")[3], ":")[1]

	slog.Debug(fmt.Sprintf("direction: %s limit: %s offset: %s", direction, limit, offset))

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
			"callback_data": fmt.Sprintf("command:friend_info;id:%s", friend.ID),
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
