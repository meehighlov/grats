package handlers

import (
	"bytes"
	"log"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

func ListBirthdaysHandler(event telegram.Event) error {
	message := event.GetMessage()
	friends, err := (&db.Friend{UserId: message.From.Id}).Filter()

	if err != nil {
		log.Println("Error fetching friends", err.Error())
		return nil
	}

	if len(friends) == 0 {
		event.Reply("Записей пока нет✨")
		return nil
	}

	var msg bytes.Buffer
	for _, friend := range friends {
		msg.WriteString(friend.Name)
		msg.WriteString(" ")
		msg.WriteString(friend.BirthDay)
		msg.WriteString("\n")
	}

	event.Reply(msg.String())

	return nil
}
