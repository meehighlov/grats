package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	models "github.com/meehighlov/grats/internal/models/call-back-data"
	"github.com/meehighlov/grats/telegram"
)

func DeleteFriendCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := models.DeleteFromRaw(event.GetCallbackQuery().Data)

	friendId := params.ID()

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error serching friend when deleting: " + err.Error())
	}

	if len(friends) == 0 {
		slog.Error("not found friend row by id: " + friendId)
		return err
	}

	friend := friends[0]

	err = friend.Delete(ctx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error delteting friend: " + err.Error())
	}

	markup := [][]map[string]string{
		{
			{
				"text": "👈к списку",
				"callback_data": models.CallList(strconv.Itoa(LIST_START_OFFSET), "<").String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, "Напоминание удалено👋", markup)

	callBackMsg := fmt.Sprintf("Напоминание для %s (%s) удалено🙌", friend.Name, friend.BirthDay)
	event.ReplyCallbackQuery(ctx, callBackMsg)

	return nil
}
