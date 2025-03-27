package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
	"gorm.io/gorm"
)

func GroupHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	invitedBy := event.GetMessage().From.Id
	if event.GetCallbackQuery().Id != "" {
		invitedBy = event.GetCallbackQuery().From.Id
	}

	// also selects supergroups
	chats, err := (&db.Chat{BotInvitedById: strconv.Itoa(invitedBy), ChatType: "%group"}).Filter(ctx, tx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(*common.NewButton("🏠 в начало", common.CallSetup().String()))
	keyboard.AppendAsStack(*common.NewAddBotToChatURLButton("➕ добавить бота в чат", config.Cfg().BotName))

	if len(chats) == 0 {
		if _, err := event.EditCalbackMessage(
			ctx,
			"После добавления в группу тут отобразится список Ваших групп, в которые я добавлен",
			*keyboard.Murkup(),
		); err != nil {
			return err
		}
		return nil
	}

	buttons := []common.Button{}

	for _, chat := range chats {
		fullInfo, err := event.GetChat(ctx, chat.ChatId)
		if err != nil {
			return err
		}
		if fullInfo != nil {
			buttons = append(buttons, *common.NewButton(fullInfo.Title, common.CallChatInfo(chat.ChatId).String()))
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

func GroupInfoHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, err := event.GetChat(ctx, params.Id)
	if err != nil {
		return err
	}

	chats, err := (&db.Chat{ChatId: params.Id}).Filter(ctx, tx)
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

func buildChatInfoMarkup(chatId string, chat *db.Chat) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	silentNotificationButtonText := "🔕 выключить звук уведомлений"
	if chat.IsAlreadySilent() {
		silentNotificationButtonText = "🔔 включить звук уведомлений"
	}

	keyboard.AppendAsStack(
		*common.NewButton("📋 список всех др в чате", common.CallChatBirthdays(chatId).String()),
		*common.NewButton("📝 изменить шаблон напоминания", common.CallEditGreetingTemplate(chatId).String()),
		*common.NewButton(silentNotificationButtonText, common.CallToggleSilentNotifications(chatId).String()),
		*common.NewButton("🗑 удалить чат", common.CallDeleteChat(chatId).String()),
		*common.NewButton("⬅️ к списку чатов", common.CallChatList().String()),
	)

	return keyboard
}

func EditGreetingTemplateHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, err := event.GetChat(ctx, params.Id)
	if err != nil {
		return err
	}

	chats, err := (&db.Chat{ChatId: params.Id}).Filter(ctx, tx)
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

func SaveGreetingTemplateHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	if len(event.GetContext().GetTexts()) == 0 {
		event.Logger.Error(
			"SaveGreetingTemplateHandler context error",
			"chatId", "is not provided on previous step",
			"userid", strconv.Itoa(event.GetMessage().From.Id),
		)
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return nil
	}

	chatId := event.GetContext().GetTexts()[0]

	newTemplate := event.GetMessage().Text

	if !strings.Contains(newTemplate, "%s") {
		event.Reply(ctx, "Шаблон должен содержать %s для подстановки имени именинника, попробуйте еще раз")
		return nil
	}

	if len(newTemplate) > 100 {
		event.Reply(ctx, "Шаблон не должен превышать 100 символов, попробуй еще раз")
		return nil
	}

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
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
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
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	chatInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Шаблон поздравления для чата `%s` обновлен!\n\nНовый шаблон:\n%s",
		chatInfo.Title,
		newTemplate)

	if _, err := event.Reply(ctx, msg, telegram.WithMarkDown()); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func DeleteChatHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
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
		*common.NewButton("⬅️ к настройкам чата", common.CallChatInfo(chatId).String()),
		*common.NewButton("🗑 удалить", common.CallConfirmDeleteChat(chatId).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Чат `%s` будет удален со всеми его напоминаниями, удаляем?", chatTitle), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ConfirmDeleteChatHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	chatId := params.Id

	chat := db.Chat{
		ChatId: chatId,
	}
	err := (&db.Friend{ChatId: chatId}).Delete(ctx, tx)
	if err != nil {
		event.Logger.Error("error deleting friends: " + err.Error())
		return err
	}

	chatInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		event.Logger.Error("error getting chat info when deleting: " + err.Error())
	}

	err = chat.Delete(ctx, tx)
	if err != nil {
		event.Logger.Error("error deleting chat: " + err.Error())
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		return err
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(*common.NewButton("⬅️к списку чатов", common.CallChatList().String()))

	chatTitle := "чат"
	if chatInfo != nil {
		chatTitle = chatInfo.Title
	}

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Чат `%s` и напоминания удалены👋", chatTitle), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ToggleSilentNotificationsHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	chatId := params.Id

	chats, err := (&db.Chat{ChatId: chatId}).Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("error getting chat: " + err.Error())
		return err
	}

	if len(chats) == 0 {
		return fmt.Errorf("chat not found")
	}

	chat := chats[0]

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

	chatInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("⚙️Настройка чата `%s`", chatInfo.Title)

	if _, err := event.EditCalbackMessage(ctx, msg, *buildChatInfoMarkup(chatId, chat).Murkup()); err != nil {
		return err
	}

	return nil
}
