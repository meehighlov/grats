package common

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

type HandlerType func(context.Context, *Event, *sql.Tx) error

func CreateRootHandler(logger *slog.Logger, handlers map[string]HandlerType) telegram.UpdateHandler {
	chatCahe := NewChatCache()
	return func(update telegram.Update, client *telegram.Client) error {
		ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
		defer cancel()

		logger.Debug("Root handler", "got update from chat", update.GetChatIdStr())

		chatContext := chatCahe.GetOrCreateChatContext(update.GetChatIdStr())
		command_ := update.Message.GetCommand()
		command := ""

		if command_ != "" {
			command = command_
			chatContext.Reset()
		} else {
			if update.CallbackQuery.Id != "" {
				params := CallbackFromString(update.CallbackQuery.Data)

				client.AnswerCallbackQuery(ctx, update.CallbackQuery.Id)

				logger.Debug("CallbackQueryHandler", "command", params.Command, "entity", params.Entity)
				logger.Info("CallbackQueryHandler", "command", params.Command, "chat id", update.GetChatIdStr())
				command = params.Command
			} else {
				command_ = chatContext.GetNextHandler()
				if command_ != "" {
					command = command_
				}
			}
		}

		if update.Message.GetChatIdStr() == config.Cfg().SupportChatId {
			command = "send_support_response"
		}

		event := newEvent(client, update, chatContext, logger)

		logger.Debug("root handler", "update", update)

		handler, found := handlers[command]
		if found {
			tx, err := db.GetDBConnection().BeginTx(ctx, nil)
			if err != nil {
				logger.Error("Root handler", "getting transaction error", err.Error())
				tx.Rollback()
			} else {
				logger.Info("Root handler", "start transaction for command", command, "chat id", update.GetChatIdStr())
				err := handler(ctx, event, tx)
				if err != nil {
					tx.Rollback()
					chatContext.Reset()
					logger.Error("Root handler", "handler error", err.Error(), "chat id", update.GetChatIdStr())
				} else {
					tx.Commit()
					logger.Info("Root handler", "transaction commited for command", command, "chat id", update.GetChatIdStr())
				}
			}
		}

		return nil
	}
}
