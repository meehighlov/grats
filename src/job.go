package src

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meehighlov/grats/db"
	"github.com/meehighlov/grats/telegram"
)

const CHECK_TIMEOUT_SEC = 10

func notify(client telegram.ApiCaller, friends []db.Friend) error {
	msg := "🔔Сегодня день рождения у %s🥳"
	for _, friend := range friends {
		msg = fmt.Sprintf(msg, friend.Name)
		_, err := client.SendMessage(friend.GetChatIdStr(), msg, false)
		if err != nil {
			log.Println("Notification not sent:", err.Error())
		}

		friend.UpdateNotifyAt()
		friend.Save()
	}

	return nil
}

func run(client telegram.ApiCaller) {
	log.Println("Starting job for checking birthdays")

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err.Error())
	}

	for {
		date := time.Now().In(location).Format("02.01.2006")

		friends, err := (&db.Friend{FilterNotifyAt: date}).Filter()

		if err != nil {
			log.Println("Error getting birthdays: " + err.Error())
		} else {
			notify(client, friends)
		}

		time.Sleep(CHECK_TIMEOUT_SEC * time.Second)
	}
}

func BirthdayNotifer(token string) error {
	client := telegram.NewClient(token)

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("Паника в фоновой задаче проверки дней рождения\n %s", r)

			log.Println(errMsg)

			reportChatId := os.Getenv("REPORT_CHAT_ID")
			_, err := client.SendMessage(reportChatId, errMsg, false)
			if err != nil {
				log.Println("panic report error:", err)
			}
		}
	}()

	run(client)

	return nil
}
