package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const HOWTO = `
1. Добавь меня в групповой чат
2. Зайди в диалог со мной и вызови /chats
3. Выбери нужный чат из списка

Напоминания будут приходить в чат в 00:00 дня рождения

Если убрать меня из чата - я безвозвратно удалю все напоминания
и не буду слать уведомления в чат 🙌
`

func GroupHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	invitedBy := event.GetMessage().From.Id
	if event.GetCallbackQuery().Id != "" {
		invitedBy = event.GetCallbackQuery().From.Id
	}

	chats, err := (&db.Chat{BotInvitedBy: strconv.Itoa(invitedBy), ChatType: "group"}).Filter(ctx, tx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(*common.NewButton("💫Как это работает?💫", common.CallChatHowto(event.GetMessage().GetChatIdStr()).String()))

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
		fullInfo, err := event.GetChat(ctx, chat.ChatId)
		if err != nil {
			return err
		}
		if fullInfo != nil {
			buttons = append(buttons, *common.NewButton(fullInfo.Title, common.CallChatInfo(chat.ChatId).String()))
		}
	}

	keyboard.AppendAsStack(buttons...)

	header := "Это чаты, в которые я добавлен✨"

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

	msg := fmt.Sprintf("Настройка чата `%s`", chatInfo.Title)

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
		*common.NewButton("добавить др в чат", common.CallAddToChat(chatId).String()),
		*common.NewButton("список всех др в чате", common.CallChatBirthdays(chatId).String()),
		*common.NewButton("⬅️к списку чатов", common.CallChatList().String()),
	)

	return keyboard
}
