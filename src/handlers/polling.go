package handlers

import (
	"github.com/meehighlov/grats/src"
	"github.com/meehighlov/grats/telegram"
)

func StartPolling(tc telegram.APICaller) error {
	cache := telegram.NewBotCache()
	updates := tc.GetUpdatesChannel()

	for update := range updates {
		if !src.IsAuthUser(update.Message.From) {
			continue
		}

		chatContext := cache.GetOrCreateChatContext(update.Message.GetChatIdStr())

		command_ := update.Message.GetCommand()
		command := ""

		if command_ != nil {
			command = *command_
			chatContext.Reset()
		} else {
			command_ = chatContext.GetCommandInProgress()
			if command_ != nil {
				command = *command_
			}
		}

		switch command {
		case "/help":
			go HelpHandler(tc, update.Message)
		case "/start":
			go StartHandler(tc, update.Message)
		case "/add":
			go AddBirthdayHandler(tc, update.Message, chatContext)
		case "/list":
			go ListBirthdaysHandler(tc, update.Message)
		default:
			continue
		}
	}

	return nil
}
