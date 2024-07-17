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

	bot.RegisterCommandHandler("/start", handlers.StartHandler)
	bot.RegisterCommandHandler("/help", handlers.HelpHandler)
	bot.RegisterCommandHandler("/list", handlers.ListBirthdaysHandler)
	bot.RegisterChatHandler("/add", handlers.AddBirthdayChatHandler())

	bot.StartPolling()
	log.Println("Polling started.")
}
