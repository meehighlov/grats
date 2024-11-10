package common

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/meehighlov/grats/telegram"
)

// use this type when there is no dialog, "one shot" handler
// for example "/start" bot command with one action
type CommandHandler func(context.Context, *Event, *sql.Tx) error

// use this type when need to build dialog with user
// FSM will invoke handlers of this type step by step
type CommandStepHandler func(context.Context, *Event, *sql.Tx) (string, error)


type Event struct {
	client  *telegram.Client
	update  telegram.Update
	context *ChatContext
	command string
	Logger  *slog.Logger
}

func newEvent(client *telegram.Client, update telegram.Update, context *ChatContext, command string, logger *slog.Logger) *Event {
	return &Event{client, update, context, command, logger}
}

func (e *Event) GetContext() *ChatContext {
	return e.context
}

func (e *Event) GetMessage() *telegram.Message {
	return &e.update.Message
}

func (e *Event) GetCallbackQuery() *telegram.CallbackQuery {
	return &e.update.CallbackQuery
}

func (e *Event) Reply(ctx context.Context, text string, opts ...telegram.SendMessageOption) (*telegram.Message, error) {
	msg, err := e.client.SendMessage(ctx, e.GetMessage().GetChatIdStr(), text, opts...)
	return msg, err
}

func (e *Event) ReplyCallbackQuery(ctx context.Context, text string, opts ...telegram.SendMessageOption) (*telegram.Message, error) {
	msg, err := e.client.SendMessage(ctx, e.GetCallbackQuery().Message.GetChatIdStr(), text, opts...)
	return msg, err
}

func (e *Event) ReplyWithKeyboard(ctx context.Context, text string, keyboard [][]map[string]string) (*telegram.Message, error) {
	msg, err := e.client.SendMessage(
		ctx,
		e.GetMessage().GetChatIdStr(),
		text,
		telegram.WithReplyMurkup(keyboard),
	)

	return msg, err
}

func (e *Event) EditCalbackMessage(ctx context.Context, text string, keyboard [][]map[string]string) (*telegram.Message, error) {
	msg, err := e.client.EditMessageText(
		ctx,
		e.update.CallbackQuery.Message.GetChatIdStr(),
		e.update.CallbackQuery.Message.GetMessageIdStr(),
		text,
		keyboard,
	)
	return msg, err
}

func (e *Event) GetChat(ctx context.Context, chatId string) (*telegram.Chat, error) {
	chat, err := e.client.GetChat(ctx, chatId)
	if chat != nil {
		return &chat.Result, err
	}
	return nil, err
}

func (e *Event) GetCommand() string {
	return e.command
}
