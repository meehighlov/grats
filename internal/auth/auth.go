package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/config"

	"github.com/meehighlov/grats/telegram"
)

func isAdmin(tgusername string) bool {
	for _, auth_user_name := range config.Cfg().AdminList() {
		if auth_user_name == tgusername {
			return true
		}
	}

	return false
}

func inAccessList(tgusername string) bool {
	hasAccess := (&db.Access{TGusername: tgusername}).IsExist(context.Background())
	return hasAccess
}

func Auth(handler telegram.CommandHandler) telegram.CommandHandler {
	return func(event telegram.Event) error {
		message := event.GetMessage()
		if isAdmin(message.From.Username) || inAccessList(message.From.Username) {
			return handler(event)
		}

		msg := fmt.Sprintf("Unauthorized access attempt by user: id=%d usernmae=%s", message.From.Id, message.From.Username)
		slog.Info(msg)

		return nil
	}
}

func Admin(handler telegram.CommandHandler) telegram.CommandHandler {
	return func(event telegram.Event) error {
		message := event.GetMessage()
		if isAdmin(message.From.Username) {
			return handler(event)
		}

		return nil
	}
}
