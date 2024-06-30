package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

const ENTER_FRIEND_NAME_STEP = 1
const ENTER_FRIEND_BIRTHDAY_STEP = 2
const SAVE_FRIEND = 3

var COMMAND = "/add"

func AddBirthdayHandler(tc telegram.APICaller, message telegram.Message, ctx telegram.ChatContext) error {
	currentStep := ctx.GetStepDone() + 1
	ctx.SetCommandInProgress(COMMAND)

	switch currentStep {
	case ENTER_FRIEND_NAME_STEP:
		enterFriendName(tc, &message)
	case ENTER_FRIEND_BIRTHDAY_STEP:
		enterBirthday(tc, &message, ctx)
	case SAVE_FRIEND:
		err := saveFriend(tc, &message, ctx)
		if err != nil {
			return nil
		}
		ctx.Reset()
		return nil
	default:
		logMsg := fmt.Sprintf("Step %d not supported for %s, resetting context", currentStep, COMMAND)
		log.Println(logMsg)
		ctx.Reset()
		return nil
	}

	ctx.SetStepDone(currentStep)

	return nil
}

func enterFriendName(tc telegram.APICaller, message *telegram.Message) error {
	msg := "Введи имя именинника✨\n\nЭто может быть 👉 имя и фамилия, никнейм и т.д."

	tc.SendMessage(message.GetChatIdStr(), msg, false)

	return nil
}

func enterBirthday(tc telegram.APICaller, message *telegram.Message, ctx telegram.ChatContext) error {
	ctx.AppendUserResponse(message.Text)

	msg := "Введи дату рождения✨\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 01.02.2003 или 01.02 "

	tc.SendMessage(message.GetChatIdStr(), msg, false)

	return nil
}

func saveFriend(tc telegram.APICaller, message *telegram.Message, ctx telegram.ChatContext) error {
	if err := validateBirthdaty(message.Text); err != nil {
		errMsg := "Не могу разобрать дату🤔\n\nВведи дату снова🙌"
		tc.SendMessage(message.GetChatIdStr(), errMsg, false)
		return err
	}

	ctx.AppendUserResponse(message.Text)
	data := ctx.GetUserResponses()
	friend := db.Friend{
		Name:     data[0],
		BirthDay: data[1],
		UserId:   message.From.Id,
		ChatId:   message.Chat.Id,
	}

	friend.Save()

	msg := fmt.Sprintf("День рождения для %s добавлен 💾", data[0])
	tc.SendMessage(message.GetChatIdStr(), msg, false)

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
