package handlers

import (
	"bytes"
	"log"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

func ListBirthdaysHandler(tc telegram.APICaller, message telegram.Message) error {
	friends, err := (&db.Friend{UserId: message.From.Id}).GetAll()

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
