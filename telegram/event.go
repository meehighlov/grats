package telegram

type Event interface {
	GetContext() ChatContext
	GetMessage() *Message
	Reply(string) *Message
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

func (e *event) Reply(text string) *Message {
	needForceReply := false
	msg, _ := e.bot.client.SendMessage(e.message.GetChatIdStr(), text, needForceReply)
	return msg
}

func (e *event) getCommand() string {
	return e.command
}

func (e *event) getStepHandlers() map[int]CommandStepHandler {
	return e.bot.chatHandlers[e.command]
}
