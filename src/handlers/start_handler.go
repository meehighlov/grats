package handlers

import (
	"fmt"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

func StartHandler(event telegram.Event) error {
	message := event.GetMessage()

	isAdmin := 0
	if message.From.IsAdmin() {
		isAdmin = 1
	}

	user := db.User{
		ID:         message.From.Id,
		Name:       message.From.FirstName,
		TGusername: message.From.Username,
		ChatId:     message.Chat.Id,
		Birthday:   "",
		IsAdmin:    isAdmin,
	}

	user.Save()

	hello := fmt.Sprintf(
		"Привет, %s 👋 Я сохраняю дни рождения и напоминаю о них🥳 \n\n /help - покажет все команды🙌",
		message.From.Username,
	)

	event.Reply(hello)

	return nil
}
