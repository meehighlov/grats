package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

const (
	MAX_CHATS_FOR_USER = 10
)

func StartHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()

	// at some point it is possible to use /command in group chat
	// so block this action
	if strings.HasSuffix(message.Chat.Type, "group") {
		return nil
	}

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

func StartFromGroupHandler(ctx context.Context, event *common.Event) error {
	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userChats, err := (&db.Chat{
			BotInvitedById: strconv.Itoa(event.GetMessage().From.Id),
			ChatType:       "%group",
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
			chat.BaseFields = db.NewBaseFields(false)
			chat.BotInvitedById = strconv.Itoa(event.GetMessage().From.Id)
			chat.GreetingTemplate = "🔔Сегодня день рождения празднует %s🥳"
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
				userChats[0].BotInvitedById,
				fmt.Sprintf("Не могу добавить новый чат, достигнут лимит (%d) подключенных групповых чатов👉👈",
					MAX_CHATS_FOR_USER))

			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
