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

const (
	ENTER_FRIEND_NAME_TO_DELETE_STEP = 1
	DELETE_FRIEND_REMINDER_STEP      = 2
	DELETE_DONE                      = -1
)

func enterFriendNameToDelete(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "Введи имя именинника, для которого нужно убрать напоминания✨"

	event.Reply(ctx, msg)

	return DELETE_FRIEND_REMINDER_STEP, nil
}

func deleteFriendReminder(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	friendName := strings.TrimSpace(event.GetMessage().Text)

	entities, err := (&db.Friend{Name: friendName}).Filter(ctx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error filtering friends while accepting name to delete: " + err.Error())
		return DONE, err
	}

	if len(entities) == 0 {
		event.Reply(ctx, "Не могу найти друга с таким именем🤔 попробуй ввести снова, учитывай верхний и нижний регистр")
		return DELETE_FRIEND_REMINDER_STEP, nil
	}

	err = (&db.Friend{Name: friendName}).Delete(ctx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error delteting friend: " + err.Error())
		return DELETE_DONE, nil
	}

	msg := fmt.Sprintf("Напоминание для %s удалено👋", friendName)

	event.Reply(ctx, msg)

	return DELETE_DONE, nil
}

func DeleteFriendChatHandler() map[int]telegram.CommandStepHandler {
	return map[int]telegram.CommandStepHandler{
		ENTER_FRIEND_NAME_TO_DELETE_STEP: enterFriendNameToDelete,
		DELETE_FRIEND_REMINDER_STEP: deleteFriendReminder,
	}
}
