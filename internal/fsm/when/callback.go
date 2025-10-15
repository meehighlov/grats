package when

import (
	"context"
	"strings"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type CallbackCondition struct {
	data string
}

// Check if callback data contains the specified data
//
// Warning:
//
// data param may match with more than one of your callbacks
// as it checks substring, so it is advised to use identification without
// repeated substrings
//
// Example:
//
// your callbacks are "item_delete", "item_delete_confirm"
// if data = "item_delete" then randon condition will be matched
//
// Tip:
//
// use idents without repeated substrings like:
// "add_item", "info_i", "1", "2", "abc", "xyz"
func CallbackDataContains(data string) *CallbackCondition {
	return &CallbackCondition{data: data}
}

func (c *CallbackCondition) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	if update.IsCallback() {
		cbd := update.CallbackQuery.Data
		return strings.Contains(cbd, c.data), nil
	}
	return false, nil
}
