package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
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

func Auth(logger *slog.Logger, handler common.HandlerType) common.HandlerType {
	return func(event common.Event) error {
		message := event.GetMessage()
		if isAdmin(message.From.Username) || inAccessList(message.From.Username) {
			return handler(event)
		}

		msg := fmt.Sprintf("Unauthorized access attempt by user: id=%d usernmae=%s", message.From.Id, message.From.Username)
		logger.Info(msg)

		return nil
	}
}

func Admin(logger *slog.Logger, handler common.HandlerType) common.HandlerType {
	return func(event common.Event) error {
		message := event.GetMessage()
		if isAdmin(message.From.Username) {
			return handler(event)
		}

		return nil
	}
}
