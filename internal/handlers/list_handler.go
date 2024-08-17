package handlers

import (
	"bytes"
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
		event.Reply(ctx, "Записей пока нет✨")
		return nil
	}

	event.ReplyWithKeyboard(
		ctx,
		listBirthdaysAsMessage(friends, LIST_LIMIT, LIST_START_OFFSET),
		buildPagiButtons(len(friends), LIST_LIMIT, LIST_START_OFFSET),
	)

	return nil
}

func listBirthdaysAsMessage(friends []db.Friend, limit, offset int) string {
	sort.Slice(friends, func(i, j int) bool { return birthdayComparator(friends, i, j) })
	var msg bytes.Buffer
	for i, friend := range friends {
		if offset != len(friends) {
			if i == limit + offset {
				break
			}
			if i < offset {
				continue
			}
		}
		msg.WriteString(friend.Name)
		msg.WriteString(" ")
		msg.WriteString(friend.BirthDay)
		if friend.IsTodayBirthday() {
			msg.WriteString(" 🥳")
		} else {
			if friend.ThisMonthAfterToday() {
				msg.WriteString(" 🕒")
			}
		}
		msg.WriteString("\n")
	}

	return msg.String()
}

func birthdayComparator(friends []db.Friend, i, j int) bool {
	bd_i := strings.Split(friends[i].BirthDay, ".")
	bd_j := strings.Split(friends[j].BirthDay, ".")

	if friends[i].IsTodayBirthday() {
		return true
	}

	if friends[j].IsTodayBirthday() {
		return false
	}

	if friends[i].ThisMonthAfterToday() {
		return true
	}

	if friends[j].ThisMonthAfterToday() {
		return false
	}

	return strings.Join([]string{bd_i[1], bd_i[0]}, ".") > strings.Join([]string{bd_j[1], bd_j[0]}, ".")
}

func buildPagiButtons(total, limit, offset int) [][]map[string]string {
	if offset == total {
		return [][]map[string]string{}
	}
	var keyBoard []map[string]string
	if offset + limit >= total {
		previousButton := map[string]string{"text": "<<", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:<<", limit, offset)}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": ">>", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:>>", limit, offset)}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": ">>", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:>>", limit, offset)}
			previousButton := map[string]string{"text": "<<", "callback_data": fmt.Sprintf("command:list;limit:%d;offset:%d;direction:<<", limit, offset)}
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

func ListBirthdaysPagination(event telegram.Event) error {
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

	if direction == ">>" {
		offset_ += LIST_PAGINATION_SHIFT
	} 
	if direction == "<<" {
		offset_ -= LIST_PAGINATION_SHIFT
	}
	if direction == "<>" {
		offset_ = len(friends)
	}

	msg := listBirthdaysAsMessage(friends, limit_, offset_)

	event.EditCalbackMessage(ctx, msg, buildPagiButtons(len(friends), limit_, offset_))

	return nil
}
