package telegram

import "context"

type Event interface {
	GetContext() ChatContext
	GetMessage() *Message
	Reply(context.Context, string) *Message
	getCommand() string
	getStepHandlers() map[int]CommandStepHandler
}

type CommandHandler func(Event) error
type CommandStepHandler func(Event) (int, error)

type event struct {
	bot     *bot
	message Message
	context ChatContext
	command string
}

func newEvent(bot *bot, message Message, context ChatContext, command string) Event {
	return &event{bot, message, context, command}
}

func (e *event) GetContext() ChatContext {
	return e.context
}

func (e *event) GetMessage() *Message {
	return &e.message
}

func (e *event) Reply(ctx context.Context, text string) *Message {
	msg, _ := e.bot.client.SendMessage(ctx, e.message.GetChatIdStr(), text)
	return msg
}

func (e *event) getCommand() string {
	return e.command
}

func (e *event) getStepHandlers() map[int]CommandStepHandler {
	return e.bot.chatHandlers[e.command]
}
