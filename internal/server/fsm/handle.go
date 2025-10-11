package fsm

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (f *FSM) Handle(ctx context.Context, update *telegram.Update) error{
	defer func() {
		if r := recover(); r != nil {
			f.logger.Error(
				"Root handler",
				"recovered from panic, error", r,
				"stack", string(debug.Stack()),
				"update", update,
			)
			f.clients.Cache.Reset(ctx, update.GetChatIdStr())

			chatId := update.GetChatIdStr()
			if chatId != "" {
				f.clients.Telegram.Reply(ctx, f.constants.ERROR_MESSAGE, update)
				return
			}

			f.logger.Error(
				"Root handler",
				"recover from panic", "chatId was empty",
				"update", update,
			)
		}
	}()

	if err := f.clients.Cache.CreateChatContext(ctx, update.GetChatIdStr()); err != nil {
		return err
	}

	if update.IsCallback() {
		f.clients.Telegram.AnswerCallbackQuery(ctx, update.CallbackQuery.Id)
	}

	command := update.GetMessage().GetCommand()
	if command != "" {
		if err := f.clients.Cache.Reset(ctx, update.GetChatIdStr()); err != nil {
			return err
		}
		return f.callNode(ctx, update, command)
	}

	handlerName, err := f.clients.Cache.GetNextHandler(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}

	return f.callNode(ctx, update, handlerName)
}

func (f *FSM) callNode(ctx context.Context, update *telegram.Update, command string) error {
	node, found := f.nodes[command]
	if !found {
		msg := "command handler not found " + command
		return errors.New(msg)
	}

	err := node.handler(ctx, update)
	if err != nil {
		nextHandler := node.getNextHandlerByError(err)
		if nextHandler == "" {
			// todo concat cache error
			if err := f.clients.Cache.Reset(ctx, update.GetChatIdStr()); err != nil {
				return err
			}
		} else {
			// todo possible cache error
			f.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), nextHandler)
		}
		return err
	}

	// todo possible cache error
	f.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), node.nextHandlerName)

	return nil
}
