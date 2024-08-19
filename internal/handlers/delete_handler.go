package handlers

import (
	"context"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

func DeleteFriendCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := strings.Split(event.GetCallbackQuery().Data, ";")

	friendId := strings.Split(params[1], ":")[1]

	baseFields := db.BaseFields{ID: friendId}
	err := (&db.Friend{BaseFields: baseFields}).Delete(ctx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error delteting friend: " + err.Error())
	}

	markup := [][]map[string]string{
		{
			{
				"text": "вернуться к списку⬅️",
				"callback_data": "command:list;limit:5;offset:0;direction:<<<",
			},
		},
	}

	event.EditCalbackMessage(ctx, "Напоминание удалено👋", markup)

	return nil
}
