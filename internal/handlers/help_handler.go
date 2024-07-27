package handlers

import (
	"context"
	"strings"

	// "github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/internal/auth"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/telegram"
)

func HelpHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить день рождения",
		"/list - список всех дней рождения",
	}

	if auth.IsAdmin(event.GetMessage().From.Username) {
		commands = append(commands, "\nАдминка🤡\n")
		commands = append(commands, "/access_list - список пользователей с доступом😏")
		commands = append(commands, "/access_grant - предоставить доступ🙈")
		commands = append(commands, "/access_revoke - отозвать доступ🤝")
	}

	msg := strings.Join(commands, "\n")

	event.Reply(ctx, msg)

	return nil
}
