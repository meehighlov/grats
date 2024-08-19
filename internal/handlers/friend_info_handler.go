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
	eventId := strings.Split(strings.Split(callbackQuery.Data, ";")[1], ":")[1]

	baseFields := db.BaseFields{ID: eventId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
		return nil
	}

	friend := friends[0]

	msgLines := []string{
		fmt.Sprintf("Имя %s", friend.Name),
		fmt.Sprintf("День рождения %s", friend.BirthDay),
		fmt.Sprintf("Напомню о нем %s", *friend.GetNotifyAt()),
		"\n",
	}

	if friend.IsTodayBirthday() {
		msgLines = append(msgLines, fmt.Sprintf("Сегодня %s празднует день рождения🥳", friend.Name))
	} else {
		if friend.IsThisMonthAfterToday() {
			msgLines = append(msgLines, fmt.Sprintf("У %s скоро день рождения🕒", friend.Name))
		}
	}

	msg := strings.Join(msgLines, "\n")

	markup := [][]map[string]string{
		{
			{
				"text": "удалить напоминание👋",
				"callback_data": fmt.Sprintf("command:delete_friend;id:%s", eventId),
			},
			{
				"text": "вернуться к списку⬅️",
				"callback_data": "command:list;limit:5;offset:0;direction:<<<",
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
