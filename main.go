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
	src.SetupFileLogging("grats.log")
	// todo load .env from basedir
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		log.Println("Error loading .env file")
		return
	}

	tc := telegram.NewClient(os.Getenv("BOT_TOKEN"))

	fmt.Println("Polling started.")
	log.Println("Polling started.")
	src.StartPolling(tc)
}
