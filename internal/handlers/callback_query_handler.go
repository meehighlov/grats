package handlers

import (
	"context"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/telegram"
)

func CallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	command := strings.Split(event.GetCallbackQuery().Data, ";")[0]

	slog.Debug("handling callback query, command: " + command)

	if command == "list" {
		ListBirthdaysCallbackQueryHandler(event)
	}
	if command == "info" {
		FriendInfoCallbackQueryHandler(event)
	}
	if command == "delete_friend" {
		DeleteFriendCallbackQueryHandler(event)
	}
	return nil
}
