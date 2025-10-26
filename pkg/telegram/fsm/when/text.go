package when

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type MessageHasTextCondition struct{}

func MessageHasText() *MessageHasTextCondition {
	return &MessageHasTextCondition{}
}

func (c *MessageHasTextCondition) Check(ctx context.Context, update *models.Update) (bool, error) {
	if update.IsCallback() {
		return false, nil
	}
	if update.GetMessage().GetCommand() != "" {
		return false, nil
	}
	return update.GetMessage().Text != "", nil
}
