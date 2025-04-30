package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/pkg/telegram"
	"gorm.io/gorm"
)

const HOWTO = `
1. Скопируйте команду, нажав кнопку ниже
2. Добавьте меня в групповой чат и отправьте туда скопированную команду

Чат должен появиться в меню "Групповые чаты"
`

func GroupHandler(ctx context.Context, event *common.Event) error {
	invitedBy := event.GetMessage().From.Id
	if event.GetCallbackQuery().Id != "" {
		invitedBy = event.GetCallbackQuery().From.Id
	}

	// also selects supergroups
	chats, err := (&db.Chat{BotInvitedById: strconv.Itoa(invitedBy), ChatType: "%group"}).Filter(ctx, nil)
	if err != nil {
		event.Reply(ctx, common.ERROR_MESSAGE)
		return err
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(common.NewButton("↩️ в начало", common.CallCommands().String()))
	keyboard.AppendAsStack(common.NewButton("💫как добавить в группу💫", common.CallChatHowto(event.GetMessage().GetChatIdStr()).String()))

	if len(chats) == 0 {
		if _, err := event.EditCalbackMessage(
			ctx,
			"Чатов пока нет🙌",
			*keyboard.Murkup(),
		); err != nil {
			return err
		}
		return nil
	}

	buttons := []*common.Button{}

	for _, chat := range chats {
		fullInfo, _ := event.GetChat(ctx, chat.ChatId)
		if fullInfo != nil {
			buttons = append(buttons, common.NewButton(fullInfo.Title, common.CallChatInfo(chat.ChatId).String()))
		}
	}

	keyboard.AppendAsStack(buttons...)

	header := "✨Это чаты, в которые я добавлен"

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			header,
			*keyboard.Murkup(),
		); err != nil {
			return err
		}
	} else {
		if _, err := event.ReplyWithKeyboard(ctx, header, *keyboard.Murkup()); err != nil {
			return err
		}
	}

	return nil
}

func GroupInfoHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, _ := event.GetChat(ctx, params.Id)
	chats, err := (&db.Chat{ChatId: params.Id}).Filter(ctx, nil)
	if err != nil {
		return err
	}

	if len(chats) == 0 {
		return fmt.Errorf("chat not found")
	}

	msg := fmt.Sprintf("⚙️Настройка чата `%s`", chatInfo.Title)

	if _, err := event.EditCalbackMessage(ctx, msg, *buildChatInfoMarkup(params.Id, chats[0]).Murkup()); err != nil {
		return err
	}

	return nil
}

func GroupHowtoHandler(ctx context.Context, event *common.Event) error {
	msg := fmt.Sprintf(
		"\nМаксимальное количество групповых чатов: %d",
		MAX_CHATS_FOR_USER,
	)

	cfg := config.Cfg()
	msg = HOWTO + msg

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		common.NewCopyButton("скопировать команду", fmt.Sprintf("/start@%s", cfg.BotName)),
		common.NewURLButton("выбрать групповой чат", fmt.Sprintf("https://t.me/%s?startgroup=true", cfg.BotName)),
		common.NewButton("⬅️ к списку чатов", common.CallChatList().String()),
	)

	if _, err := event.EditCalbackMessage(
		ctx,
		msg,
		*keyboard.Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func buildChatInfoMarkup(chatId string, chat *db.Chat) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	silentNotificationButtonText := "🔕 выключить звук уведомлений"
	if chat.IsAlreadySilent() {
		silentNotificationButtonText = "🔔 включить звук уведомлений"
	}

	keyboard.AppendAsStack(
		common.NewButton("📋 список всех др в чате", common.CallList(strconv.Itoa(LIST_START_OFFSET), ">", chatId, "friend").String()),
		common.NewButton("📝 изменить шаблон напоминания", common.CallEditGreetingTemplate(chatId).String()),
		common.NewButton(silentNotificationButtonText, common.CallToggleSilentNotifications(chatId).String()),
		common.NewButton("🗑 удалить чат", common.CallDeleteChat(chatId).String()),
		common.NewButton("⬅️ к списку чатов", common.CallChatList().String()),
	)

	return keyboard
}

func EditGreetingTemplateHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, err := event.GetChat(ctx, params.Id)
	if err != nil {
		return err
	}

	chats, err := (&db.Chat{ChatId: params.Id}).Filter(ctx, nil)
	if err != nil {
		return err
	}

	currentTemplate := ""
	if len(chats) > 0 && chats[0].GreetingTemplate != "" {
		currentTemplate = chats[0].GreetingTemplate
	}

	header := fmt.Sprintf("Шаблон для чата `%s`:\n\n%s\n",
		chatInfo.Title,
		currentTemplate,
	)

	msg := strings.Join([]string{
		header,
		"Напишите новый шаблон в ответ на это сообщение, используйте %s для подстановки имени именинника",
	}, "\n")

	event.ReplyCallbackQuery(ctx, msg, telegram.WithMarkDown())

	event.GetContext().AppendText(params.Id)
	event.SetNextHandler("save_greeting_template")

	return nil
}

func SaveGreetingTemplateHandler(ctx context.Context, event *common.Event) error {
	chatId := event.GetContext().GetTexts()[0]
	newTemplate := event.GetMessage().Text

	if len(event.GetContext().GetTexts()) == 0 {
		event.Logger.Error(
			"SaveGreetingTemplateHandler context error",
			"chatId", "is not provided on previous step",
			"userid", strconv.Itoa(event.GetMessage().From.Id),
		)
		event.Reply(ctx, common.ERROR_MESSAGE)
		return nil
	}

	if !strings.Contains(newTemplate, "%s") {
		event.Reply(ctx, "Шаблон должен содержать %s для подстановки имени именинника, попробуйте еще раз")
		return nil
	}

	if len(newTemplate) > 100 {
		event.Reply(ctx, "Шаблон не должен превышать 100 символов, попробуй еще раз")
		return nil
	}

	done := false

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chats, err := (&db.Chat{ChatId: chatId, BotInvitedById: strconv.Itoa(event.GetMessage().From.Id)}).Filter(ctx, tx)
		if err != nil {
			return err
		}

		if len(chats) == 0 {
			event.Logger.Error(
				"SaveGreetingTemplateHandler chats not found",
				"chatId", chatId,
				"userid", strconv.Itoa(event.GetMessage().From.Id),
			)
			event.Reply(ctx, common.ERROR_MESSAGE)
			return nil
		}

		chat := chats[0]
		chat.GreetingTemplate = newTemplate

		err = chat.Save(ctx, tx)
		if err != nil {
			event.Logger.Error(
				"SaveGreetingTemplateHandler db error",
				"chatId", chatId,
				"userid", strconv.Itoa(event.GetMessage().From.Id),
				"error", err,
			)
			event.Reply(ctx, common.ERROR_MESSAGE)
			return err
		}

		done = true

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		event.SetNextHandler("")
		chatInfo, _ := event.GetChat(ctx, chatId)
		msg := fmt.Sprintf("Шаблон поздравления для чата `%s` обновлен!\n\nНовый шаблон:\n%s",
			chatInfo.Title,
			newTemplate)

		if _, err := event.Reply(ctx, msg, telegram.WithMarkDown()); err != nil {
			return err
		}
	}
	return nil
}

func DeleteChatHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	chatId := params.Id

	chatInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		event.Logger.Error("error getting chat info when deleting: " + err.Error())
		return err
	}

	chatTitle := "чат"
	if chatInfo != nil {
		chatTitle = chatInfo.Title
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		common.NewButton("⬅️ к настройкам чата", common.CallChatInfo(chatId).String()),
		common.NewButton("🗑 удалить", common.CallConfirmDeleteChat(chatId).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Чат `%s` будет удален со всеми его напоминаниями, удаляем?", chatTitle), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ConfirmDeleteChatHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	chatId := params.Id
	keyboard := common.NewInlineKeyboard()

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chat := db.Chat{
			ChatId: chatId,
		}
		err := (&db.Friend{ChatId: chatId}).Delete(ctx, tx)
		if err != nil {
			event.Logger.Error("error deleting friends: " + err.Error())
			return err
		}

		err = chat.Delete(ctx, tx)
		if err != nil {
			event.Logger.Error("error deleting chat: " + err.Error())
			if _, err := event.ReplyCallbackQuery(ctx, common.ERROR_MESSAGE); err != nil {
				return err
			}
			return err
		}

		keyboard.AppendAsStack(common.NewButton("⬅️к списку чатов", common.CallChatList().String()))

		return nil
	})

	if err != nil {
		return err
	}

	chatInfo, _ := event.GetChat(ctx, chatId)
	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Чат `%s` и напоминания удалены👋", chatInfo.Title), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ToggleSilentNotificationsHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	chatId := params.Id

	done := false

	var chat *db.Chat

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chats, err := (&db.Chat{ChatId: chatId}).Filter(ctx, tx)
		if err != nil {
			event.Logger.Error("error getting chat: " + err.Error())
			return err
		}

		if len(chats) == 0 {
			return fmt.Errorf("chat not found")
		}

		chat = chats[0]

		if chat.IsAlreadySilent() {
			chat.EnableSoundNotifications()
		} else {
			chat.DisableSoundNotifications()
		}

		err = chat.Save(ctx, tx)
		if err != nil {
			event.Logger.Error("error saving chat: " + err.Error())
			return err
		}

		done = true

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		chatInfo, err := event.GetChat(ctx, chatId)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("⚙️Настройка чата `%s`", chatInfo.Title)

		if _, err := event.EditCalbackMessage(ctx, msg, *buildChatInfoMarkup(chatId, chat).Murkup()); err != nil {
			return err
		}
	}

	return nil
}
