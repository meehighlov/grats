package telegram

func (bot *bot) StartPolling() error {
	updates := bot.client.GetUpdatesChannel()

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
