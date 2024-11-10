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

	markup := [][]map[string]string{}
	howtoButton := map[string]string{
		"text":          "💫Как это работает?💫",
		"callback_data": common.CallChatHowto(event.GetMessage().GetChatIdStr()).String(),
	}
	markup = append(markup, []map[string]string{howtoButton})

	if len(chats) == 0 {
		if _, err := event.ReplyWithKeyboard(
			ctx,
			"Чатов пока нет🙌",
			markup,
		); err != nil {
			return err
		}
		return nil
	}

	buttons := []map[string]string{}

	for _, chat := range chats {
		fullInfo, err := event.GetChat(ctx, chat.ChatId)
		if err != nil {
			return err
		}
		if fullInfo != nil {
			button := map[string]string{
				"text":          fullInfo.Title,
				"callback_data": common.CallChatInfo(chat.ChatId).String(),
			}
			buttons = append(buttons, button)
		}
	}

	for _, button := range buttons {
		markup = append(markup, []map[string]string{button})
	}

	header := "Это чаты, в которые я добавлен✨"

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			header,
			markup,
		); err != nil {
			return err
		}
	} else {
		if _, err := event.ReplyWithKeyboard(ctx, header, markup); err != nil {
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

	if _, err := event.EditCalbackMessage(ctx, msg, buildChatInfoMarkup(params.Id)); err != nil {
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

func buildChatInfoMarkup(chatId string) [][]map[string]string {
	addBirthDayButton := map[string]string{
		"text":          "добавить др в чат",
		"callback_data": common.CallAddToChat(chatId).String(),
	}
	listButton := map[string]string{
		"text":          "список всех др в чате",
		"callback_data": common.CallChatBirthdays(chatId).String(),
	}
	backButton := map[string]string{
		"text":          "к списку чатов",
		"callback_data": common.CallChatList().String(),
	}

	markup := [][]map[string]string{}
	markup = append(markup, []map[string]string{addBirthDayButton})
	markup = append(markup, []map[string]string{listButton})
	markup = append(markup, []map[string]string{backButton})

	return markup
}
