package handlers

import (
	"context"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

func HelpHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить день рождения",
		"/list - список всех дней рождения",
	}

	message := event.GetMessage()
	u, err := (&db.User{TGId: message.From.Id, TGusername: message.From.Username}).Filter(ctx)

	if err != nil {
		slog.Error("Error filtering users when building help command: " + err.Error())
	} else {

		if len(u) == 1 && u[0].HasAdminAccess() {
			commands = append(commands, "\nАдминка🤡\n")
			commands = append(commands, "/access_list - список пользователей с доступом😏")
			commands = append(commands, "/access_grant - предоставить доступ🙈")
			commands = append(commands, "/access_revoke - отозвать доступ🤝")
		}
	}

	msg := strings.Join(commands, "\n")

	event.Reply(ctx, msg)

	return nil
}
