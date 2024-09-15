package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

func FriendInfoCallbackQueryHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
		return nil
	}

	friend := friends[0]

	// todo take from db
	friendTimezone := "мск"

	emoji, zodiacName := friend.GetZodiacSign()

	msgLines := []string{
		fmt.Sprintf("✨ %s", friend.Name),
		fmt.Sprintf("🗓 %s", friend.BirthDay),
		fmt.Sprintf("%s %s", emoji, zodiacName),
		fmt.Sprintf("🔔 Напомню %s в полночь (по %s)", *friend.GetNotifyAt(), friendTimezone),
	}

	if friend.IsTodayBirthday() {
		msgLines = append(msgLines, fmt.Sprintf("🥳 Сегодня %s празднует день рождения", friend.Name))
	} else {
		if friend.IsThisMonthAfterToday() {
			msgLines = append(msgLines, fmt.Sprintf("🕒 У %s скоро день рождения", friend.Name))
		}
	}

	msg := strings.Join(msgLines, "\n\n")

	offset := params.Pagination.Offset

	markup := [][]map[string]string{
		{
			{
				"text": "👈к списку",
				"callback_data": common.CallList(offset, "<").String(),
			},
		},
		{
			{
				"text": "удалить👋",
				"callback_data": common.CallDelete(params.Id).String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
