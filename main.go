package main

import (
	"github.com/joho/godotenv"
	"github.com/meehighlov/grats/telegram"
	"github.com/meehighlov/grats/src"
	"log"
	"os"
)

func main() {
	src.SetupFileLogging("grats.log")
	// todo load .env from basedir
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	tc := telegram.NewClient(os.Getenv("BOT_TOKEN"))

	log.Println("Polling started.")
	src.StartPolling(tc)
}
