package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

func FriendInfoCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	callbackQuery := event.GetCallbackQuery()
	params := strings.Split(callbackQuery.Data, ";")
	friendId := params[1]

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
		return nil
	}

	friend := friends[0]

	// todo take from db
	friendTimezone := "мск"

	msgLines := []string{
		fmt.Sprintf("🟢 %s", friend.Name),
		fmt.Sprintf("🟢 %s", friend.BirthDay),
		fmt.Sprintf("🟢 Напомню о нем %s в полночь, часовой пояс - %s", *friend.GetNotifyAt(), friendTimezone),
	}

	if friend.IsTodayBirthday() {
		msgLines = append(msgLines, fmt.Sprintf("🥳 Сегодня %s празднует день рождения", friend.Name))
	} else {
		if friend.IsThisMonthAfterToday() {
			msgLines = append(msgLines, fmt.Sprintf("🕒 У %s скоро день рождения", friend.Name))
		}
	}

	msg := strings.Join(msgLines, "\n\n")

	offset := params[2]

	markup := [][]map[string]string{
		{
			{
				"text": "👈к списку",
				"callback_data": fmt.Sprintf("list;%s;<", offset),
			},
		},
		{
			{
				"text": "удалить👋",
				"callback_data": fmt.Sprintf("delete_friend;%s;", friendId),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
