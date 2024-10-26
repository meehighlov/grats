package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func DeleteFriendCallbackQueryHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	friendId := params.Id

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error serching friend when deleting: " + err.Error())
		return err
	}

	if len(friends) == 0 {
		slog.Error("not found friend row by id: " + friendId)
		return err
	}

	friend := friends[0]

	err = friend.Delete(ctx, tx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error delteting friend: " + err.Error())
		return err
	}

	markup := [][]map[string]string{
		{
			{
				"text":          "👈к списку",
				"callback_data": common.CallList(strconv.Itoa(LIST_START_OFFSET), "<", params.BoundChat).String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, "Напоминание удалено👋", markup)

	callBackMsg := fmt.Sprintf("Напоминание для %s (%s) удалено🙌", friend.Name, friend.BirthDay)
	event.ReplyCallbackQuery(ctx, callBackMsg)

	return nil
}
