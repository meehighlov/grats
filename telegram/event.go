package telegram

import "context"

type Event interface {
	GetContext() ChatContext
	GetMessage() *Message
	GetCallbackQuery() *CallbackQuery
	Reply(context.Context, string, ...sendMessageOption) *Message
	ReplyCallbackQuery(context.Context, string, ...sendMessageOption) *Message
	ReplyWithKeyboard(context.Context, string, [][]map[string]string) *Message
	EditCalbackMessage(context.Context, string, [][]map[string]string) *Message
	AnswerCallbackQuery(context.Context) bool
	getCommand() string
	getStepHandlers() map[int]CommandStepHandler
}

type CommandHandler func(Event) error
type CommandStepHandler func(Event) (int, error)

type event struct {
	bot     *bot
	update  Update
	context ChatContext
	command string
}

func newEvent(bot *bot, update Update, context ChatContext, command string) Event {
	return &event{bot, update, context, command}
}

func (e *event) GetContext() ChatContext {
	return e.context
}

func (e *event) GetMessage() *Message {
	return &e.update.Message
}

func (e *event) GetCallbackQuery() *CallbackQuery {
	return &e.update.CallbackQuery
}

func (e *event) Reply(ctx context.Context, text string, opts ...sendMessageOption) *Message {
	msg, _ := e.bot.client.SendMessage(ctx, e.GetMessage().GetChatIdStr(), text, opts...)
	return msg
}

func (e *event) ReplyCallbackQuery(ctx context.Context, text string, opts ...sendMessageOption) *Message {
	msg, _ := e.bot.client.SendMessage(ctx, e.GetCallbackQuery().Message.GetChatIdStr(), text, opts...)
	return msg
}

func (e *event) ReplyWithKeyboard(ctx context.Context, text string, keyboard [][]map[string]string) *Message {
	msg, _ := e.bot.client.SendMessage(
		ctx,
		e.GetMessage().GetChatIdStr(),
		text,
		WithReplyMurkup(keyboard),
	)

	return msg
}

func (e *event) EditCalbackMessage(ctx context.Context, text string, keyboard [][]map[string]string) *Message {
	msg, _ := e.bot.client.EditMessageText(
		ctx,
		e.update.CallbackQuery.Message.GetChatIdStr(),
		e.update.CallbackQuery.Message.GetMessageIdStr(),
		text,
		keyboard,
	)
	return msg
}

func (e *event) AnswerCallbackQuery(ctx context.Context) bool {
	err := e.bot.client.AnswerCallbackQuery(ctx, e.update.CallbackQuery.Id)
	if err != nil {
		return false
	} else {
		return false
	}
}

func (e *event) getCommand() string {
	return e.command
}

func (e *event) getStepHandlers() map[int]CommandStepHandler {
	return e.bot.chatHandlers[e.command]
}
