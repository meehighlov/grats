package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

const GRANT_ACCESS_ENTRYPOINT = 1
const SAVE_TG_USERNAME = 2

const REVOKE_ACCESS_ENTRYPOINT = 1
const UPDATE_ACCESS_INFO = 2

func AccessListHandler(event telegram.Event) error {
	accessList, err := (&db.Access{}).All()

	if err != nil {
		log.Println("Error fetching access list", err.Error())
		event.Reply(err.Error())
		return nil
	}

	if len(*accessList) == 0 {
		event.Reply("В таблице доступов нет записей✨")
		return nil
	}

	var msg bytes.Buffer
	for _, access := range *accessList {
		msg.WriteString(access.GetTGUserName())
		msg.WriteString("\n")
	}

	event.Reply(msg.String())

	return nil
}

func grantAccess(event telegram.Event) (int, error) {
	msg := "Кому предоставить доступ? Введи имя пользователя тг😘"

	event.Reply(msg)

	return SAVE_TG_USERNAME, nil
}

func saveAccess(event telegram.Event) (int, error) {
	tgusername := event.GetMessage().Text
	tgusername = strings.Replace(tgusername, "@", "", 1)

	err := (&db.Access{BaseFields: db.NewBaseFields(), TGusername: tgusername}).Save()

	if err != nil {
		event.Reply(err.Error())
		return SAVE_TG_USERNAME, nil
	}

	msg := fmt.Sprintf("Доступ для %s предоставлен, пусть пробует зайти💋", tgusername)

	event.Reply(msg)

	return telegram.STEPS_DONE, nil
}

func revokeAccess(event telegram.Event) (int, error) {
	msg := "У кого отбираем доступ?😡"

	event.Reply(msg)

	return UPDATE_ACCESS_INFO, nil
}

func updateAccessInfo(event telegram.Event) (int, error) {
	tgusername := strings.Replace(event.GetMessage().Text, "@", "", 1)
	err := (&db.Access{TGusername: tgusername}).Delete()

	if err != nil {
		event.Reply(err.Error())
		return UPDATE_ACCESS_INFO, nil
	}

	msg := fmt.Sprintf("Доступ для %s закрыт🖐", event.GetMessage().Text)

	event.Reply(msg)

	return telegram.STEPS_DONE, nil
}

func GrantAccessChatHandler() map[int]telegram.CommandStepHandler {
	handlers := make(map[int]telegram.CommandStepHandler)

	handlers[GRANT_ACCESS_ENTRYPOINT] = grantAccess
	handlers[SAVE_TG_USERNAME] = saveAccess

	return handlers
}

func RevokeAccessChatHandler() map[int]telegram.CommandStepHandler {
	handlers := make(map[int]telegram.CommandStepHandler)

	handlers[REVOKE_ACCESS_ENTRYPOINT] = revokeAccess
	handlers[UPDATE_ACCESS_INFO] = updateAccessInfo

	return handlers
}
