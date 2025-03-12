package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

const HOWTO = `
1. Добавь меня в групповой чат
2. Вызови /start в групповом чате
3. Если в ответ в чат придет сообщение "Всем привет👋",
   значит настройка прошла успешно

После шага 3 чат отобразится в меню "Групповые чаты"
Уведомления будут приходить в групповой чат
`

func GroupHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	invitedBy := event.GetMessage().From.Id
	if event.GetCallbackQuery().Id != "" {
		invitedBy = event.GetCallbackQuery().From.Id
	}

	// also selects supergroups
	chats, err := (&db.Chat{BotInvitedBy: strconv.Itoa(invitedBy), ChatType: "%group"}).Filter(ctx, tx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(*common.NewButton("🏠 в начало", common.CallSetup().String()))
	keyboard.AppendAsStack(*common.NewButton("💫Как настроить отправку напоминаний в чаты?💫", common.CallChatHowto(event.GetMessage().GetChatIdStr()).String()))

	if len(chats) == 0 {
		if _, err := event.ReplyWithKeyboard(
			ctx,
			"Чатов пока нет🙌",
			*keyboard.Murkup(),
		); err != nil {
			return err
		}
		return nil
	}

	buttons := []common.Button{}

	for _, chat := range chats {
		fullInfo, err := event.GetChat(ctx, chat.TGChatId)
		if err != nil {
			return err
		}
		if fullInfo != nil {
			buttons = append(buttons, *common.NewButton(fullInfo.Title, common.CallChatInfo(chat.TGChatId).String()))
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

func GroupInfoHandler(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, err := event.GetChat(ctx, params.Id)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("⚙️Настройка чата `%s`", chatInfo.Title)

	if _, err := event.EditCalbackMessage(ctx, msg, *buildChatInfoMarkup(params.Id).Murkup()); err != nil {
		return err
	}

	return nil
}

func GroupHowtoHandler(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	if _, err := event.ReplyCallbackQuery(ctx, HOWTO); err != nil {
		return err
	}

	return nil
}

func buildChatInfoMarkup(chatId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		*common.NewButton("📋 список всех др в чате", common.CallChatBirthdays(chatId).String()),
		*common.NewButton("🔔 изменить шаблон напоминания", common.CallEditGreetingTemplate(chatId).String()),
		*common.NewButton("🗑 удалить чат", common.CallDeleteChat(chatId).String()),
		*common.NewButton("⬅️ к списку чатов", common.CallChatList().String()),
	)

	return keyboard
}

func EditGreetingTemplateHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatInfo, err := event.GetChat(ctx, params.Id)
	if err != nil {
		return err
	}

	chats, err := (&db.Chat{TGChatId: params.Id}).Filter(ctx, tx)
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
		"Пришли новый шаблон в ответ на это сообщение, используй %s для подстановки имени именинника",
	}, "\n")

	event.ReplyCallbackQuery(ctx, msg, telegram.WithMarkDown())

	event.GetContext().AppendText(params.Id)
	event.SetNextHandler("save_greeting_template")

	return nil
}

func SaveGreetingTemplateHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	if len(event.GetContext().GetTexts()) == 0 {
		event.Logger.Error(
			"SaveGreetingTemplateHandler context error",
			"chatId", "is not provided on previous step",
			"userid", strconv.Itoa(event.GetMessage().From.Id),
		)
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return nil
	}

	tgChatId := event.GetContext().GetTexts()[0]

	newTemplate := event.GetMessage().Text

	if !strings.Contains(newTemplate, "%s") {
		event.Reply(ctx, "Шаблон должен содержать %s для подстановки имени именинника, попробуй еще раз")
		return nil
	}

	if len(newTemplate) > 100 {
		event.Reply(ctx, "Шаблон не должен превышать 100 символов, попробуй еще раз")
		return nil
	}

	chats, err := (&db.Chat{TGChatId: tgChatId, BotInvitedBy: strconv.Itoa(event.GetMessage().From.Id)}).Filter(ctx, tx)
	if err != nil {
		return err
	}

	if len(chats) == 0 {
		event.Logger.Error(
			"SaveGreetingTemplateHandler chats not found",
			"tgChatId", tgChatId,
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
			"tgChatId", tgChatId,
			"userid", strconv.Itoa(event.GetMessage().From.Id),
			"error", err,
		)
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	chatInfo, err := event.GetChat(ctx, tgChatId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Шаблон поздравления для чата `%s` успешно обновлен!\n\nНовый шаблон:\n%s",
		chatInfo.Title,
		newTemplate)

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		*common.NewButton("⬅️к настройкам чата", common.CallChatInfo(tgChatId).String()),
	)

	if _, err := event.ReplyWithKeyboard(ctx, msg, *keyboard.Murkup()); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func DeleteChatHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
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

func ConfirmDeleteChatHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	tgChatId := params.Id

	chat := db.Chat{
		TGChatId: tgChatId,
	}
	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("error getting chat: " + err.Error())
		return err
	}

	if len(chats) == 0 {
		event.Logger.Error("chat not found: " + tgChatId)
		return err
	}

	err = (&db.Friend{ChatId: chats[0].ID}).Delete(ctx, tx)
	if err != nil {
		event.Logger.Error("error deleting friends: " + err.Error())
		return err
	}

	chatInfo, err := event.GetChat(ctx, tgChatId)
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
