package src

import (
	"fmt"
	"log"
	"strings"
	"bytes"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

func StartHandler(tc telegram.APICaller, message telegram.Message) error {

	isAdmin := 0
	if message.From.IsAdmin() {
		isAdmin = 1
	}

	user := db.User{
		ID:         message.From.Id,
		Name:       message.From.FirstName,
		TGusername: message.From.Username,
		ChatId:     message.Chat.Id,
		Birthday:   "",
		IsAdmin:    isAdmin,
	}

	user.Save()

	hello := fmt.Sprintf(
		"Привет, %s 👋 Я сохраняю дни рождения и напоминаю о них🥳 \n\n /help - покажет все команды🙌",
		message.From.Username,
	)

	tc.SendMessage(message.GetChatIdStr(), hello, false)

	return nil
}

func AddBirthdayHandler(tc telegram.APICaller, message telegram.Message) error {
	// truncate command, so start from index 1
	data := strings.Split(message.Text, " ")[1:]

	if len(data) != 2 {
		tc.SendMessage(message.GetChatIdStr(), "Не могу сделать такую запись🤔 Формат такой: /add Имя dd.mm.yyyy", false)

		return nil
	}

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

func ListBirthdaysHandler(tc telegram.APICaller, message telegram.Message) error {
	friends, err := (&db.Friend{}).All(message.From.Id)

	if err != nil {
		log.Println("Error fetching friends", err.Error())
		return nil
	}

	if len(friends) == 0 {
		tc.SendMessage(message.GetChatIdStr(), "Записей пока нет✨", false)
		return nil
	}

	var msg bytes.Buffer
	for _, friend := range friends {
		msg.WriteString(friend.Name)
		msg.WriteString(" ")
		msg.WriteString(friend.BirthDay)
		msg.WriteString("\n")
	}

	tc.SendMessage(message.GetChatIdStr(), msg.String(), false)

	return nil
}

func DeleteBirtdayHandler(tc telegram.APICaller, message telegram.Message) error {
	// deletes birthday from db for user
	// sends success to chat

	return nil
}

func HelpHandler(tc telegram.APICaller, message telegram.Message) error {
	commands := []string{
		"Ниже - список команд с примерами использования🙌",
		"\n",
		"/add Имя 01.02.2003 - добавить запись",
		"/list - список всех записей",
	}

	msg := strings.Join(commands, "\n")

	tc.SendMessage(message.GetChatIdStr(), msg, false)

	return nil
}

func StartPolling(tc telegram.APICaller) error {

	updates := tc.GetUpdatesChannel()

	for update := range updates {
		if !IsAuthUser(update.Message.From) {
			continue
		}

		command := update.Message.GetCommand()

		// todo use map

		switch command {
		case "/help":
			go HelpHandler(tc, update.Message)
		case "/start":
			go StartHandler(tc, update.Message)
		case "/add":
			go AddBirthdayHandler(tc, update.Message)
		case "/list":
			go ListBirthdaysHandler(tc, update.Message)
		case "/delete":
			go DeleteBirtdayHandler(tc, update.Message)
		default:
			continue
		}
	}

	return nil
}
