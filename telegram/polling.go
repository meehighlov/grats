package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func IsAuthUser(user User) bool {
	// temporary func
	// todo: move to middleware
	for _, auth_user_name := range strings.Split(os.Getenv("AUTH_USERS"), ",") {
		if auth_user_name == user.Username {
			return true
		}
	}

	return false
}

func (bot *bot) StartPolling() error {
	updates := bot.client.GetUpdatesChannel()

	for update := range updates {
		if !IsAuthUser(update.Message.From) {
			continue
		}

		chatContext := bot.getOrCreateChatContext(update.Message.GetChatIdStr())

		command_ := update.Message.GetCommand()
		command := ""

		if command_ != "" {
			command = command_
			chatContext.reset()
		} else {
			command_ = chatContext.getCommandInProgress()
			if command_ != "" {
				command = command_
			}
		}

		event := newEvent(bot, update.Message, chatContext)

		commandHandler, found := bot.commandHandlers[command]

		if found {
			go commandHandler(event)
		} else {
			stepHandlers, found := bot.chatHandlers[command]
			if found {
				go invokeStepHandler(command, event, stepHandlers)
			}
		}
	}

	return nil
}

func invokeStepHandler(command string, event Event, handlers map[int]CommandStepHandler) error {
	ctx := event.GetContext()
	stepTODO := ctx.getStepTODO()
	ctx.setCommandInProgress(command)

	nextStep := -1

	stepHandler, found := handlers[stepTODO]

	if !found {
		logMsg := fmt.Sprintf("Step %d not supported for %s, resetting context", stepTODO, command)
		log.Println(logMsg)
		ctx.reset()
		return nil
	}

	nextStep, _ = stepHandler(event)

	if nextStep == -1 {
		ctx.reset()
		return nil
	}

	ctx.setStepTODO(nextStep)

	return nil
}
