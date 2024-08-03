package telegram

import "context"

func (bot *bot) StartPolling() error {
	withCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates := bot.client.GetUpdatesChannel(withCancel)

	for update := range updates {
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

		event := newEvent(bot, update.Message, chatContext, command)

		commandHandler, found := bot.commandHandlers[command]

		if found {
			go commandHandler(event)
		}
	}

	return nil
}
