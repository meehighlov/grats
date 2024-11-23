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

func GrantAccess(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	msg := "Кому предоставить доступ? Введи имя пользователя тг😘"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("access_save_tg_username")

	return nil
}

func SaveAccess(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	tgusername := event.GetMessage().Text
	tgusername = strings.Replace(tgusername, "@", "", 1)

	err := (&db.Access{BaseFields: db.NewBaseFields(), TGusername: tgusername}).Save(ctx, tx)

	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return err
		}
		return err
	}

	msg := fmt.Sprintf("Доступ для %s предоставлен, пусть пробует зайти💋", tgusername)

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func RevokeAccess(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	msg := "У кого отбираем доступ?😡"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("access_update")

	return nil
}

func UpdateAccessInfo(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	tgusername := strings.Replace(event.GetMessage().Text, "@", "", 1)
	err := (&db.Access{TGusername: tgusername}).Delete(ctx, tx)

	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return err
		}
		event.SetNextHandler("access_update")
		return err
	}

	msg := fmt.Sprintf("Доступ для %s закрыт🖐", event.GetMessage().Text)

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}
