package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	models "github.com/meehighlov/grats/internal/models/call-back-data"
	"github.com/meehighlov/grats/telegram"
)

func FriendInfoCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	callbackQuery := event.GetCallbackQuery()

	info := models.InfoFromRaw(callbackQuery.Data)

	baseFields := db.BaseFields{ID: info.ID()}
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

	_, offset, _ := info.Pagination()

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
				"callback_data": fmt.Sprintf("delete;%s", info.ID()),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
