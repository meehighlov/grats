package handlers

import (
	"strings"

	"github.com/meehighlov/grats/src"
	"github.com/meehighlov/grats/telegram"
)

func HelpHandler(event telegram.Event) error {
	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить день рождения",
		"/list - список всех дней рождения",
	}

	if src.IsAdmin(event.GetMessage().From.Username) {
		commands = append(commands, "\nАдминка🤡\n")
		commands = append(commands, "/access_list - список пользователей с доступом😏")
		commands = append(commands, "/access_grant - предоставить доступ🙈")
		commands = append(commands, "/access_revoke - отозвать доступ🤝")
	}

	msg := strings.Join(commands, "\n")

	event.Reply(msg)

	return nil
}
