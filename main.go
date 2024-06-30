package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/meehighlov/grats/src"
	"github.com/meehighlov/grats/telegram"
	"log"
	"os"
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

	tc := telegram.NewClient(botToken)

	fmt.Println("Polling started.")
	log.Println("Polling started.")
	src.StartPolling(tc)
}
