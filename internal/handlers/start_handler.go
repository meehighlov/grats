package handlers

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/common"
)

func StartHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()

	err := RegisterOrUpdateUser(ctx, event)
	if err != nil {
		event.Logger.Error("start error registering user", "chatId", message.GetChatIdStr(), "error", err.Error())
		return err
	}

	username := message.From.Username
	if username == "" {
		username = message.From.FirstName
		if username == "" {
			username = "друг"
		}
	}

	hello := fmt.Sprintf(
		("Привет, %s👋" +
			"\n\n" +
			"/commands - покажет все мои команды"),
		message.From.Username,
	)

	if _, err := event.Reply(ctx, hello); err != nil {
		return err
	}

	return nil
}
