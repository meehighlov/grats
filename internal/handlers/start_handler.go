package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func StartHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	message := event.GetMessage()

	isAdmin := 0
	if message.From.IsAdmin() {
		isAdmin = 1
	}

	user := db.User{
		BaseFields: db.NewBaseFields(),
		Name:       message.From.FirstName,
		TGusername: message.From.Username,
		TGId:       message.From.Id,
		ChatId:     message.Chat.Id,
		Birthday:   "",
		IsAdmin:    isAdmin,
	}

	err := user.Save(ctx, tx)
	if err != nil {
		return err
	}

	chat := db.Chat{
		BaseFields:   db.NewBaseFields(),
		ChatType:     "private",
		ChatId:       event.GetMessage().GetChatIdStr(),
		BotInvitedBy: event.GetMessage().From.Id,
	}

	err = chat.Save(ctx, tx)
	if err != nil {
		return err
	}

	hello := fmt.Sprintf(
		"Привет, %s 👋 Я напоминаю о днях рождения🥳",
		message.From.Username,
	)

	event.Reply(ctx, hello)

	return nil
}
