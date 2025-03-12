package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func StartHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()

	isAdmin := 0
	if message.From.IsAdmin() {
		isAdmin = 1
	}

	user := db.User{
		BaseFields: db.NewBaseFields(),
		Name:       message.From.FirstName,
		TGusername: message.From.Username,
		TGId:       strconv.Itoa(message.From.Id),
		ChatId:     strconv.Itoa(message.Chat.Id),
		Birthday:   "",
		IsAdmin:    isAdmin,
	}

	err := user.Save(ctx, tx)
	if err != nil {
		return err
	}

	_, err = db.GetOrCreateChatByTGChatId(
		ctx,
		tx,
		event.GetMessage().GetChatIdStr(),
		"private",
		strconv.Itoa(event.GetMessage().From.Id),
	)
	if err != nil {
		return err
	}

	hello := fmt.Sprintf(
		("Привет, %s👋 Меня зовут grats" +
			"\n" +
			"Я напоминаю о днях рождения🥳" +
			"\n\n" +
			"Команда /setup покажет все мои команды"),
		message.From.Username,
	)

	if _, err := event.Reply(ctx, hello); err != nil {
		return err
	}

	return nil
}

func StartFromGroupHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	chatType := event.GetMessage().Chat.Type

	_, err := db.GetOrCreateChatByTGChatId(
		ctx,
		tx,
		event.GetMessage().GetChatIdStr(),
		chatType,
		strconv.Itoa(event.GetMessage().From.Id),
	)
	if err != nil {
		event.Reply(ctx, "Что-то пошло не так🙃 Попробуй еще раз позже👉👈")
		return err
	}

	event.Reply(ctx, "Всем привет👋")

	return nil
}
