package handlers

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

func ListBirthdaysHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()
	friends, err := (&db.Friend{UserId: message.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching friends" + err.Error())
		return nil
	}

	if len(friends) == 0 {
		event.Reply(ctx, "Записей пока нет✨")
		return nil
	}

	var msg bytes.Buffer
	for _, friend := range friends {
		msg.WriteString(friend.Name)
		msg.WriteString(" ")
		msg.WriteString(friend.BirthDay)
		msg.WriteString("\n")
	}

	event.Reply(ctx, msg.String())

	return nil
}
