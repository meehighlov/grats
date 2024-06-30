package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

const ENTER_FRIEND_NAME_STEP = 1
const ENTER_FRIEND_BIRTHDAY_STEP = 2
const SAVE_FRIEND_STEP = 3
const DONE = -1

func EnterFriendName(event telegram.Event) (int, error) {
	msg := "Введи имя именинника✨\n\nнапример 👉 Райан Гослинг"

	event.Reply(msg)

	return ENTER_FRIEND_BIRTHDAY_STEP, nil
}

func EnterBirthday(event telegram.Event) (int, error) {
	event.GetContext().AppendUserResponse(event.GetMessage().Text)

	msg := "Введи дату рождения✨\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 12.11.1980 или 12.11"

	event.Reply(msg)

	return SAVE_FRIEND_STEP, nil
}

func SaveFriend(event telegram.Event) (int, error) {
	message := event.GetMessage()
	ctx := event.GetContext()

	if err := validateBirthdaty(message.Text); err != nil {
		errMsg := "Дата не попадает под формат🤔\n\nвведи дату снова🙌"
		event.Reply(errMsg)
		return SAVE_FRIEND_STEP, err
	}

	ctx.AppendUserResponse(message.Text)
	data := ctx.GetUserResponses()

	friend := db.Friend{
		BaseFields: db.NewBaseFields(),
		Name:     data[0],
		BirthDay: data[1],
		UserId:   message.From.Id,
		ChatId:   message.Chat.Id,
	}

	friend.RenewNotifayAt()

	friend.Save()

	msg := fmt.Sprintf("День рождения для %s добавлен 💾", data[0])
	event.Reply(msg)

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

func AddBirthdayChatHandler() map[int]telegram.CommandStepHandler {
	handlers := make(map[int]telegram.CommandStepHandler)

	handlers[ENTER_FRIEND_NAME_STEP] = EnterFriendName
	handlers[ENTER_FRIEND_BIRTHDAY_STEP] = EnterBirthday
	handlers[SAVE_FRIEND_STEP] = SaveFriend

	return handlers
}
