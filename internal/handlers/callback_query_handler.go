package handlers

import (
	"context"
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/models"
	"github.com/meehighlov/grats/telegram"
)

func CallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	command := models.CallbackFromString(event.GetCallbackQuery().Data).Command

	slog.Debug("handling callback query, command: " + command)

	if command == "list" {
		ListBirthdaysCallbackQueryHandler(event)
	}
	if command == "info" {
		FriendInfoCallbackQueryHandler(event)
	}
	if command == "delete" {
		DeleteFriendCallbackQueryHandler(event)
	}
	return nil
}
