package models

type ChatContext struct {
	ChatId        string   `json:"chatId"`
	UserResponses []string `json:"userResponses"`
	StateStatus   string   `json:"stateStatus"`
}

func (ctx *ChatContext) AppendText(userResponse string) error {
	ctx.UserResponses = append(ctx.UserResponses, userResponse)
	return nil
}

func (ctx *ChatContext) GetTexts() []string {
	return ctx.UserResponses
}

func (ctx *ChatContext) GetStateStatus() string {
	return ctx.StateStatus
}

func (ctx *ChatContext) Reset() error {
	ctx.UserResponses = []string{}
	return nil
}
