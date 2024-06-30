package src

import (
	"log"
	"fmt"
	"os"
)

func SetupFileLogging(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error configuring logging", err.Error())
		log.Panic(err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	return logFile
}
