package common

import (
	"context"
	"log/slog"

	"github.com/meehighlov/grats/telegram"
)

type Event struct {
	client  *telegram.Client
	update  telegram.Update
	context *ChatContext
	Logger  *slog.Logger
}

func newEvent(client *telegram.Client, update telegram.Update, context *ChatContext, logger *slog.Logger) *Event {
	return &Event{client, update, context, logger}
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

func (e *Event) GetNextHandler() string {
	return e.GetContext().GetNextHandler()
}

func (e *Event) SetNextHandler(nextHandler string) string {
	return e.GetContext().SetNextHandler(nextHandler)
}
