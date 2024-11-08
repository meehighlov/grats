package common

import (
	"context"
	"database/sql"

	"github.com/meehighlov/grats/telegram"
)

// todo return error in all methods
type Event interface {
	GetContext() *ChatContext
	GetMessage() *telegram.Message
	GetCallbackQuery() *telegram.CallbackQuery
	Reply(context.Context, string, ...telegram.SendMessageOption) *telegram.Message
	ReplyCallbackQuery(context.Context, string, ...telegram.SendMessageOption) *telegram.Message
	ReplyWithKeyboard(context.Context, string, [][]map[string]string) *telegram.Message
	EditCalbackMessage(context.Context, string, [][]map[string]string) *telegram.Message
	GetChat(context.Context, string) *telegram.Chat
	GetCommand() string
}

type CommandHandler func(context.Context, Event, *sql.Tx) error
type CommandStepHandler func(context.Context, Event, *sql.Tx) (string, error)

type event struct {
	client  telegram.ApiCaller
	update  telegram.Update
	context *ChatContext
	command string
}

func newEvent(client telegram.ApiCaller, update telegram.Update, context *ChatContext, command string) Event {
	return &event{client, update, context, command}
}

func (e *event) GetContext() *ChatContext {
	return e.context
}

func (e *event) GetMessage() *telegram.Message {
	return &e.update.Message
}

func (e *event) GetCallbackQuery() *telegram.CallbackQuery {
	return &e.update.CallbackQuery
}

func (e *event) Reply(ctx context.Context, text string, opts ...telegram.SendMessageOption) *telegram.Message {
	msg, _ := e.client.SendMessage(ctx, e.GetMessage().GetChatIdStr(), text, opts...)
	return msg
}

func (e *event) ReplyCallbackQuery(ctx context.Context, text string, opts ...telegram.SendMessageOption) *telegram.Message {
	msg, _ := e.client.SendMessage(ctx, e.GetCallbackQuery().Message.GetChatIdStr(), text, opts...)
	return msg
}

func (e *event) ReplyWithKeyboard(ctx context.Context, text string, keyboard [][]map[string]string) *telegram.Message {
	msg, _ := e.client.SendMessage(
		ctx,
		e.GetMessage().GetChatIdStr(),
		text,
		telegram.WithReplyMurkup(keyboard),
	)

	return msg
}

func (e *event) EditCalbackMessage(ctx context.Context, text string, keyboard [][]map[string]string) *telegram.Message {
	msg, _ := e.client.EditMessageText(
		ctx,
		e.update.CallbackQuery.Message.GetChatIdStr(),
		e.update.CallbackQuery.Message.GetMessageIdStr(),
		text,
		keyboard,
	)
	return msg
}

func (e *event) GetChat(ctx context.Context, chatId string) *telegram.Chat {
	chat, _ := e.client.GetChat(ctx, chatId)
	if chat != nil {
		return &chat.Result
	}
	return nil
}

func (e *event) GetCommand() string {
	return e.command
}
