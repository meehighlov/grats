package handlers

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

// Существующий хендлер для автоматической регистрации группового чата
func GroupChatRegisterHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	chat := db.Chat{
		BaseFields:       db.NewBaseFields(),
		ChatType:         "group",
		BotInvitedBy:     strconv.Itoa(event.GetMessage().From.Id),
		ChatId:           event.GetMessage().GetChatIdStr(),
		GreetingTemplate: "🔔Сегодня день рождения у %s🥳",
	}

	message := event.GetMessage()

	if message.LeftChatMember.Username == config.Cfg().BotName {
		// todo check bot name by GetMe tg method
		err := chat.Delete(ctx, tx)
		if err != nil {
			return err
		}
		err = (&db.Friend{ChatId: strconv.Itoa(message.Chat.Id)}).Delete(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	}

	if message.NewChatMembers != nil {
		for _, member := range message.NewChatMembers {
			// todo check bot name by GetMe tg method
			if member.Username == config.Cfg().BotName {
				err := chat.Save(ctx, tx)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Новый хендлер для ручной регистрации чата через команду
func AddChatHandler(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	msg := "Введи ID чата для добавления✨\n\nнапример 👉 -1001234567890"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_chat_save")

	return nil
}

// Хендлер для сохранения информации о чате
func SaveChatHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	chatId := event.GetMessage().Text

	// Проверка, что чат существует
	chatInfo, err := event.GetChat(ctx, chatId)
	if err != nil || chatInfo == nil {
		if _, err := event.Reply(ctx, "Не могу найти такой чат. Убедись, что бот добавлен в чат и ID указан верно."); err != nil {
			return err
		}
		event.SetNextHandler("add_chat_save")
		return nil
	}

	// Создаем новую запись чата
	chat := db.Chat{
		BaseFields:       db.NewBaseFields(),
		ChatType:         "group",
		BotInvitedBy:     strconv.Itoa(event.GetMessage().From.Id),
		ChatId:           chatId,
		GreetingTemplate: "🔔Сегодня день рождения у %s🥳",
	}

	// Проверяем, существует ли уже такой чат
	existingChats, err := (&db.Chat{ChatId: chatId}).Filter(ctx, tx)
	if err != nil {
		if _, err := event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		return err
	}

	if len(existingChats) > 0 {
		if _, err := event.Reply(ctx, "Этот чат уже добавлен в систему!"); err != nil {
			return err
		}
		event.SetNextHandler("")
		return nil
	}

	// Сохраняем чат
	err = chat.Save(ctx, tx)
	if err != nil {
		if _, err := event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		return err
	}

	// Создаем клавиатуру для навигации
	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		*common.NewButton("⬅️ к списку чатов", common.CallChatList().String()),
		*common.NewButton("настройки чата", common.CallChatInfo(chatId).String()),
	)

	msg := "Чат успешно добавлен! Теперь ты можешь добавлять дни рождения для этого чата."

	if _, err := event.ReplyWithKeyboard(ctx, msg, *keyboard.Murkup()); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}
