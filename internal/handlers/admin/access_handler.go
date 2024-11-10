package admin

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	GRANT_ACCESS_ENTRYPOINT  = "1"
	SAVE_TG_USERNAME         = "2"
	REVOKE_ACCESS_ENTRYPOINT = "1"
	UPDATE_ACCESS_INFO       = "2"
)

func AdminCommandListHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	commands := []string{
		"/access_list - список пользователей с доступом😏",
		"/access_grant - предоставить доступ🙈",
		"/access_revoke - отозвать доступ🤝",
	}

	msg := strings.Join(commands, "\n")

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	return nil
}

func AccessListHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	accessList, err := (&db.Access{}).All(ctx, tx)

	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return err
		}
		return err
	}

	if len(*accessList) == 0 {
		if _, err := event.Reply(ctx, "В таблице доступов нет записей✨"); err != nil {
			return err
		}
		return err
	}

	var msg bytes.Buffer
	for _, access := range *accessList {
		msg.WriteString(access.GetTGUserName())
		msg.WriteString("\n")
	}

	if _, err := event.Reply(ctx, msg.String()); err != nil {
		return err
	}

	return nil
}

func grantAccess(ctx context.Context, event *common.Event, _ *sql.Tx) (string, error) {
	msg := "Кому предоставить доступ? Введи имя пользователя тг😘"

	if _, err := event.Reply(ctx, msg); err != nil {
		return common.STEPS_DONE, err
	}

	return SAVE_TG_USERNAME, nil
}

func saveAccess(ctx context.Context, event *common.Event, tx *sql.Tx) (string, error) {
	tgusername := event.GetMessage().Text
	tgusername = strings.Replace(tgusername, "@", "", 1)

	err := (&db.Access{BaseFields: db.NewBaseFields(), TGusername: tgusername}).Save(ctx, tx)

	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return common.STEPS_DONE, err
		}
		return SAVE_TG_USERNAME, err
	}

	msg := fmt.Sprintf("Доступ для %s предоставлен, пусть пробует зайти💋", tgusername)

	if _, err := event.Reply(ctx, msg); err != nil {
		return common.STEPS_DONE, err
	}

	return common.STEPS_DONE, nil
}

func revokeAccess(ctx context.Context, event *common.Event, _ *sql.Tx) (string, error) {
	msg := "У кого отбираем доступ?😡"

	if _, err := event.Reply(ctx, msg); err != nil {
		return common.STEPS_DONE, err
	}

	return UPDATE_ACCESS_INFO, nil
}

func updateAccessInfo(ctx context.Context, event *common.Event, tx *sql.Tx) (string, error) {
	tgusername := strings.Replace(event.GetMessage().Text, "@", "", 1)
	err := (&db.Access{TGusername: tgusername}).Delete(ctx, tx)

	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return common.STEPS_DONE, err
		}
		return UPDATE_ACCESS_INFO, err
	}

	msg := fmt.Sprintf("Доступ для %s закрыт🖐", event.GetMessage().Text)

	if _, err := event.Reply(ctx, msg); err != nil {
		return common.STEPS_DONE, err
	}

	return common.STEPS_DONE, nil
}

func GrantAccessChatHandler() map[string]common.CommandStepHandler {
	handlers := make(map[string]common.CommandStepHandler)

	handlers[GRANT_ACCESS_ENTRYPOINT] = grantAccess
	handlers[SAVE_TG_USERNAME] = saveAccess

	return handlers
}

func RevokeAccessChatHandler() map[string]common.CommandStepHandler {
	handlers := make(map[string]common.CommandStepHandler)

	handlers[REVOKE_ACCESS_ENTRYPOINT] = revokeAccess
	handlers[UPDATE_ACCESS_INFO] = updateAccessInfo

	return handlers
}
