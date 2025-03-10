package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	FRIEND_NAME_MAX_LEN = 50
	EMPTY_CHAT_ID       = "empty"
)

func AddToChatHandler(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	msg := "Введи имя именинника✨\n\nнапример 👉 Райан Гослинг"

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}
	event.GetContext().AppendText(common.CallbackFromString(event.GetCallbackQuery().Data).Id)

	event.SetNextHandler("add_enter_bd")

	return nil
}

func EnterBirthday(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	friendName := strings.TrimSpace(event.GetMessage().Text)

	if len(friendName) > FRIEND_NAME_MAX_LEN {
		if _, err := event.Reply(ctx, fmt.Sprintf("Имя не должно превышать %d символов", FRIEND_NAME_MAX_LEN)); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_bd")
		return nil
	}

	chatId := event.GetContext().GetTexts()[0]

	entities, err := (&db.Friend{Name: friendName, ChatId: chatId}).Filter(ctx, tx)
	if err != nil {
		if _, err := event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error filtering friends while accepting name to save: " + err.Error())
		return err
	}

	if len(entities) != 0 {
		if _, err := event.Reply(ctx, "Такое имя уже есть😅 попробуй снова, учитывай верхний и нижний регистр букв"); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_bd")
		return nil
	}

	event.GetContext().AppendText(friendName)

	msg := "Введи дату рождения✨\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 12.11.1980 или 12.11"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_save_friend")

	return nil
}

func SaveFriend(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()
	chatContext := event.GetContext()

	if err := validateBirthdaty(message.Text); err != nil {
		errMsg := "Дата не попадает под формат🤔\n\nвведи дату снова🙌"
		if _, err := event.Reply(ctx, errMsg); err != nil {
			return err
		}
		event.SetNextHandler("add_save_friend")
		return nil
	}

	chatContext.AppendText(message.Text)
	data := chatContext.GetTexts()
	chatid, name, bd := data[0], data[1], data[2]

	friend := db.Friend{
		BaseFields: db.NewBaseFields(),
		Name:       name,
		BirthDay:   bd,
		UserId:     strconv.Itoa(message.From.Id),
		ChatId:     chatid,
	}

	friend.RenewNotifayAt()

	err := friend.Save(ctx, tx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("День рождения для %s добавлен 💾\n\nНапомню тебе о нем %s🔔", name, *friend.GetNotifyAt())

	if strings.Contains(chatid, "-") {
		chatTitle := "чат"
		chatFullInfo, err := event.GetChat(ctx, chatid)
		if err != nil {
			return err
		}
		if chatFullInfo.Id != 0 {
			chatTitle = fmt.Sprintf("чат %s", chatFullInfo.Title)
		}

		msg = fmt.Sprintf("День рождения для %s добавлен в %s 💾\n\nПришлю напоминание в чат %s🔔", name, chatTitle, *friend.GetNotifyAt())
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		msg,
		*buildNavigationMarkup(chatid).Murkup(),
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func validateBirthdaty(birtday string) error {
	month := "01"
	day := "02"
	format_wo_year := fmt.Sprintf("%s.%s", day, month)
	format_with_year := fmt.Sprintf("%s.%s.2006", day, month)

	format := ""

	parts := strings.Split(birtday, ".")
	if len(parts) == 3 {
		format = format_with_year
	} else {
		format = format_wo_year
	}

	_, err := time.Parse(format, birtday)

	if err != nil {
		return err
	}

	return nil
}

func buildNavigationMarkup(chatId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(*common.NewButton("добавить еще", common.CallAddToChat(chatId).String()), *common.NewButton("список др", common.CallChatBirthdays(chatId).String()))

	return keyboard
}
