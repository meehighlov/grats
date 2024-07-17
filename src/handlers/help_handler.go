package handlers

import (
	"strings"

	"github.com/meehighlov/grats/telegram"
)

func HelpHandler(event telegram.Event) error {
	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить день рождения",
		"/list - список всех дней рождения",
	}

	msg := strings.Join(commands, "\n")

	event.Reply(msg)

	return nil
}
