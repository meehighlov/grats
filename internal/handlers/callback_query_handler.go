package handlers

import (
	"context"
	"strings"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/telegram"
)

func CallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	command := strings.Split(strings.Split(event.GetCallbackQuery().Data, ";")[0], ":")[1]
	if command == "list" {
		ListBirthdaysPagination(event)
	}
	return nil
}
