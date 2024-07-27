package handlers

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

const (
	GRANT_ACCESS_ENTRYPOINT  = 1
	SAVE_TG_USERNAME         = 2
	REVOKE_ACCESS_ENTRYPOINT = 1
	UPDATE_ACCESS_INFO       = 2
)

func AccessListHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	accessList, err := (&db.Access{}).All(ctx)

	if err != nil {
		event.Reply(ctx, err.Error())
		return nil
	}

	if len(*accessList) == 0 {
		event.Reply(ctx, "В таблице доступов нет записей✨")
		return nil
	}

	var msg bytes.Buffer
	for _, access := range *accessList {
		msg.WriteString(access.GetTGUserName())
		msg.WriteString("\n")
	}

	event.Reply(ctx, msg.String())

	return nil
}

func grantAccess(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "Кому предоставить доступ? Введи имя пользователя тг😘"

	event.Reply(ctx, msg)

	return SAVE_TG_USERNAME, nil
}

func saveAccess(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	tgusername := event.GetMessage().Text
	tgusername = strings.Replace(tgusername, "@", "", 1)

	err := (&db.Access{BaseFields: db.NewBaseFields(), TGusername: tgusername}).Save(ctx)

	if err != nil {
		event.Reply(ctx, err.Error())
		return SAVE_TG_USERNAME, nil
	}

	msg := fmt.Sprintf("Доступ для %s предоставлен, пусть пробует зайти💋", tgusername)

	event.Reply(ctx, msg)

	return telegram.STEPS_DONE, nil
}

func revokeAccess(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "У кого отбираем доступ?😡"

	event.Reply(ctx, msg)

	return UPDATE_ACCESS_INFO, nil
}

func updateAccessInfo(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	tgusername := strings.Replace(event.GetMessage().Text, "@", "", 1)
	err := (&db.Access{TGusername: tgusername}).Delete(ctx)

	if err != nil {
		event.Reply(ctx, err.Error())
		return UPDATE_ACCESS_INFO, nil
	}

	msg := fmt.Sprintf("Доступ для %s закрыт🖐", event.GetMessage().Text)

	event.Reply(ctx, msg)

	return telegram.STEPS_DONE, nil
}

func GrantAccessChatHandler() map[int]telegram.CommandStepHandler {
	handlers := make(map[int]telegram.CommandStepHandler)

	handlers[GRANT_ACCESS_ENTRYPOINT] = grantAccess
	handlers[SAVE_TG_USERNAME] = saveAccess

	return handlers
}

func RevokeAccessChatHandler() map[int]telegram.CommandStepHandler {
	handlers := make(map[int]telegram.CommandStepHandler)

	handlers[REVOKE_ACCESS_ENTRYPOINT] = revokeAccess
	handlers[UPDATE_ACCESS_INFO] = updateAccessInfo

	return handlers
}
