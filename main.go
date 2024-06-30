package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/meehighlov/grats/src"
	"github.com/meehighlov/grats/src/handlers"
	"github.com/meehighlov/grats/telegram"
)

func main() {
	logFile := src.SetupFileLogging("grats.log")
	defer logFile.Close()

	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory. Exiting.")
		return
	}

	err = godotenv.Load(curDir + "/.env")
	if err != nil {
		fmt.Println("Error loading .env file. Exiting.")

		return
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Println("BOT_TOKEN is not set. Exiting.")
		return
	}

	bot := telegram.NewBot(botToken)

	go src.BirthdayNotifer(botToken)

	bot.RegisterCommandHandler("/start", src.Auth(handlers.StartHandler))
	bot.RegisterCommandHandler("/help", src.Auth(handlers.HelpHandler))
	bot.RegisterCommandHandler("/list", src.Auth(handlers.ListBirthdaysHandler))
	bot.RegisterCommandHandler("/add", src.Auth(telegram.FSM(handlers.AddBirthdayChatHandler())))
	bot.RegisterCommandHandler("/access_list", src.Admin(handlers.AccessListHandler))
	bot.RegisterCommandHandler("/access_grant", src.Admin(telegram.FSM(handlers.GrantAccessChatHandler())))
	bot.RegisterCommandHandler("/access_revoke", src.Admin(telegram.FSM(handlers.RevokeAccessChatHandler())))

	bot.StartPolling()
	log.Println("Polling started.")
}
