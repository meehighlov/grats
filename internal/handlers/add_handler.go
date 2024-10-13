package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

const (
	ENTER_FRIEND_NAME_STEP     = "1"
	ENTER_FRIEND_BIRTHDAY_STEP = "2"
	SAVE_FRIEND_STEP           = "3"
	DONE                       = "done"

	FRIEND_NAME_MAX_LEN = 50
)

func enterFriendName(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "Введи имя именинника✨\n\nнапример 👉 Райан Гослинг"

	event.Reply(ctx, msg)

	return ENTER_FRIEND_BIRTHDAY_STEP, nil
}

func enterBirthday(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	friendName := strings.TrimSpace(event.GetMessage().Text)

	if len(friendName) > FRIEND_NAME_MAX_LEN {
		event.Reply(ctx, fmt.Sprintf("Имя не должно превышать %d символов", FRIEND_NAME_MAX_LEN))
		return ENTER_FRIEND_BIRTHDAY_STEP, nil
	}

	entities, err := (&db.Friend{Name: friendName}).Filter(ctx)
	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		slog.Error("error filtering friends while accepting name to save: " + err.Error())
		return DONE, err
	}

	if len(entities) != 0 {
		event.Reply(ctx,"Такое имя уже есть😅 попробуй снова, учитывай верхний и нижний регистр букв")
		return ENTER_FRIEND_BIRTHDAY_STEP, nil
	}

	event.GetContext().AppendText(friendName)

	msg := "Введи дату рождения✨\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 12.11.1980 или 12.11"

	event.Reply(ctx, msg)

	return SAVE_FRIEND_STEP, nil
}

func saveFriend(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()
	chatContext := event.GetContext()

	if err := validateBirthdaty(message.Text); err != nil {
		errMsg := "Дата не попадает под формат🤔\n\nвведи дату снова🙌"
		event.Reply(ctx, errMsg)
		return SAVE_FRIEND_STEP, err
	}

	chatContext.AppendText(message.Text)
	data := chatContext.GetTexts()

	friend := db.Friend{
		BaseFields: db.NewBaseFields(),
		Name:       data[0],
		BirthDay:   data[1],
		UserId:     message.From.Id,
		ChatId:     message.Chat.Id,
	}

	friend.RenewNotifayAt()

	friend.Save(context.Background())

	msg := fmt.Sprintf("День рождения для %s добавлен 💾\n\nНапомню тебе о нем %s🔔", data[0], *friend.GetNotifyAt())
	event.Reply(ctx, msg)

	return DONE, nil
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

func AddBirthdayChatHandler() map[string]common.CommandStepHandler {
	handlers := make(map[string]common.CommandStepHandler)

	handlers[ENTER_FRIEND_NAME_STEP] = enterFriendName
	handlers[ENTER_FRIEND_BIRTHDAY_STEP] = enterBirthday
	handlers[SAVE_FRIEND_STEP] = saveFriend

	return handlers
}
