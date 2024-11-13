package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func DeleteFriendCallbackQueryHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	friendId := params.Id

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error serching friend when deleting: " + err.Error())
		return err
	}

	if len(friends) == 0 {
		event.Logger.Error("not found friend row by id: " + friendId)
		return err
	}

	friend := friends[0]

	err = friend.Delete(ctx, tx)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error delteting friend: " + err.Error())
		return err
	}

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(*common.NewButton("⬅️", common.CallList(strconv.Itoa(LIST_START_OFFSET), "<", params.BoundChat).String()))

	if _, err := event.EditCalbackMessage(ctx, "Напоминание удалено👋", *keyboard.Murkup()); err != nil {
		return err
	}

	callBackMsg := fmt.Sprintf("Напоминание для %s (%s) удалено🙌", friend.Name, friend.BirthDay)
	if _, err := event.ReplyCallbackQuery(ctx, callBackMsg); err != nil {
		return err
	}

	return nil
}
