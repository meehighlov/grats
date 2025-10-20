package when

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type MessageHasTextCondition struct{}

func MessageHasText() *MessageHasTextCondition {
	return &MessageHasTextCondition{}
}

func (c *MessageHasTextCondition) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	if update.IsCallback() {
		return false, nil
	}
	if update.GetMessage().GetCommand() != "" {
		return false, nil
	}
	return update.GetMessage().Text != "", nil
}
