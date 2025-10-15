package when

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type UpdateHasOnlyTextCondition struct{}

func UpdateHasOnlyText() *UpdateHasOnlyTextCondition {
	return &UpdateHasOnlyTextCondition{}
}

func (c *UpdateHasOnlyTextCondition) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	if update.IsCallback() {
		return false, nil
	}
	if update.GetMessage().GetCommand() != "" {
		return false, nil
	}
	return true, nil
}
