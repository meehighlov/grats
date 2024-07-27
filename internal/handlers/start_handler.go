package handlers

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

func StartHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

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

	user.Save(ctx)

	hello := fmt.Sprintf(
		"Привет, %s 👋 Я сохраняю дни рождения и напоминаю о них🥳 \n\n /help - покажет все команды🙌",
		message.From.Username,
	)

	event.Reply(ctx, hello)

	return nil
}
