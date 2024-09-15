package common

import (
	"context"
	"log/slog"

	"github.com/meehighlov/grats/telegram"
)

type HandlerType func(Event) error


func CreateRootHandler(logger *slog.Logger, chatCahe *ChatCache, handlers map[string]HandlerType) telegram.UpdateHandler {
	return func(update telegram.Update, client telegram.ApiCaller) error {
		chatContext := chatCahe.GetOrCreateChatContext(update.Message.GetChatIdStr())
		command_ := update.Message.GetCommand()
		command := ""

		if command_ != "" {
			command = command_
			chatContext.Reset()
		} else {
			if update.CallbackQuery.Id != "" {
				params := CallbackFromString(update.CallbackQuery.Data)

				logger.Debug("CallbackQueryHandler", "command", params.Command, "entity", params.Entity)
				command = params.Command

				client.AnswerCallbackQuery(context.Background(), update.CallbackQuery.Id)
			} else {
				command_ = chatContext.GetCommandInProgress()
				if command_ != "" {
					command = command_
				}
			}
		}

		event := newEvent(client, update, chatContext, command)

		handler, found := handlers[command]
		if found {
			handler(event)
		}

		return nil
	}
}
