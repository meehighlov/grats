package telegram

type Event interface {
	GetContext() ChatContext
	GetMessage() *Message
	Reply(string) *Message
}

type CommandHandler func(Event) error
type CommandStepHandler func(Event) (int, error)

type event struct {
	bot     *bot
	message Message
	context ChatContext
}

func newEvent(bot *bot, message Message, context ChatContext) Event {
	return &event{bot, message, context}
}

func (e *event) GetContext() ChatContext {
	return e.context
}

func (e *event) GetMessage() *Message {
	return &e.message
}

func (e *event) Reply(text string) *Message {
	needForceReply := false
	msg := e.bot.client.SendMessage(e.message.GetChatIdStr(), text, needForceReply)
	return msg
}
