package handlers

import (
	"strings"

	"github.com/meehighlov/grats/telegram"
)

func HelpHandler(tc telegram.APICaller, message telegram.Message) error {
	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить день рождения",
		"/list - список всех дней рождения",
	}

	msg := strings.Join(commands, "\n")

	tc.SendMessage(message.GetChatIdStr(), msg, false)

	return nil
}
