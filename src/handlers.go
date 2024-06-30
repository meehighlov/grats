package src

import (
	"fmt"
	"log"
	"strings"

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

	tc.SendMessage(message.GetChatIdStr(), "Привет!", false)

	return nil
}

func AddBirthdayHandler(tc telegram.APICaller, message telegram.Message) error {
	// todo validation

	// truncate command, so start from index 1
	data := strings.Split(message.Text, " ")[1:]

	friend := db.Friend{
		Name:     data[0],
		BirthDay: data[1],
		UserId:   message.From.Id,
		ChatId:   message.Chat.Id,
	}

	friend.Save()

	msg := fmt.Sprintf("День рождения для %s добавлен", data[0])

	tc.SendMessage(message.GetChatIdStr(), msg, false)

	return nil
}

func ListBirthdaysHandler(tc telegram.APICaller, message telegram.Message) error {
	// todo add output info

	friends, err := db.Friend{}.All(message.From.Id)

	if err != nil {
		log.Println("Error fetching friends", err.Error())
	}

	msg := []string{}
	for _, friend := range friends {

		msg = append(msg, friend.Name)

	}

	tc.SendMessage(message.GetChatIdStr(), strings.Join(msg, "\n"), false)

	return nil
}

func DeleteBirtdayHandler(tc telegram.APICaller, message telegram.Message) error {
	// deletes birthday from db for user
	// sends success to chat

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

		// todo add /help command hanlder

		switch command {
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
