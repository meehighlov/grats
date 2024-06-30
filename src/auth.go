package src

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

func IsAdmin(tgusername string) bool {
	for _, auth_user_name := range strings.Split(os.Getenv("ADMINS"), ",") {
		if auth_user_name == tgusername {
			return true
		}
	}

	return false
}

func inAccessList(tgusername string) bool {
	hasAccess := (&db.Access{TGusername: tgusername}).IsExist()
	return hasAccess
}

func Auth(handler telegram.CommandHandler) telegram.CommandHandler {
	return func(event telegram.Event) error {
		message := event.GetMessage()
		if IsAdmin(message.From.Username) || inAccessList(message.From.Username) {
			return handler(event)
		}

		msg := fmt.Sprintf("Unauthorized access attempt by user: id=%d usernmae=%s", message.From.Id, message.From.Username)
		log.Println(msg)

		return nil
	}
}

func Admin(handler telegram.CommandHandler) telegram.CommandHandler {
	return func(event telegram.Event) error {
		message := event.GetMessage()
		if IsAdmin(message.From.Username) {
			return handler(event)
		}

		return nil
	}
}
