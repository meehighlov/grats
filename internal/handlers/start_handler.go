package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	MAX_CHATS_FOR_USER = 10
)

func StartHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()

	// at some point it is possible to use /command in group chat
	// so block this action
	if strings.HasSuffix(message.Chat.Type, "group") {
		return nil
	}

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

	chat := db.Chat{
		BaseFields:   db.NewBaseFields(),
		ChatType:     "private",
		ChatId:       event.GetMessage().GetChatIdStr(),
		BotInvitedBy: strconv.Itoa(event.GetMessage().From.Id),
	}

	err = chat.Save(ctx, tx)
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
	userChats, err := (&db.Chat{
		BotInvitedBy: strconv.Itoa(event.GetMessage().From.Id),
		ChatType:     "%group",
	}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error(
			"StartFromGroupHandler",
			"chat", event.GetMessage().GetChatIdStr(),
			"userId", event.GetMessage().From.Id,
			"error", err.Error(),
		)
		return err
	}

	chatType := event.GetMessage().Chat.Type
	chat := db.Chat{
		ChatId: event.GetMessage().GetChatIdStr(),
	}

	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		event.Logger.Error(
			"StartFromGroupHandler",
			"chat", chat.ChatId,
			"userId", event.GetMessage().From.Id,
			"error", err.Error(),
		)
		return err
	}

	if len(chats) == 0 && len(userChats) < MAX_CHATS_FOR_USER {
		chat.BaseFields = db.NewBaseFields()
		chat.BotInvitedBy = strconv.Itoa(event.GetMessage().From.Id)
		chat.GreetingTemplate = "🔔Сегодня день рождения у %s🥳"
		chat.ChatType = chatType

		err := chat.Save(ctx, tx)
		if err != nil {
			event.Logger.Error(
				"StartFromGroupHandler",
				"chat", chat.ChatId,
				"userId", event.GetMessage().From.Id,
				"error", err.Error(),
			)
			event.Reply(ctx, "Что-то пошло не так🙃 Попробуйте еще раз👉👈")
			return nil
		}

		event.Reply(ctx, "Всем привет👋")
		return nil
	}

	if len(chats) == 0 && len(userChats) >= MAX_CHATS_FOR_USER {
		event.Logger.Info(
			"StartFromGroupHandler",
			"chat", event.GetMessage().GetChatIdStr(),
			"userId", event.GetMessage().From.Id,
			"error", "user reached chats limits",
		)
		event.ReplyToUser(
			ctx,
			userChats[0].BotInvitedBy,
			fmt.Sprintf("Не могу добавить новый чат, достигнут лимит (%d) подключенных групповых чатов👉👈",
				MAX_CHATS_FOR_USER))

		return nil
	}

	event.Reply(ctx, "Всем привет👋")

	return nil
}
