package telegram

import (
	"context"
	"strings"
)

func MessageHasText() Condition {
	return func(ctx context.Context, scope *Scope) (bool, error) {
		if scope.Update().IsCallback() {
			return false, nil
		}
		if scope.Update().GetMessage().GetCommand() != "" {
			return false, nil
		}
		return scope.Update().GetMessage().Text != "", nil
	}
}

func Command(command string) Condition {
	return func(ctx context.Context, scope *Scope) (bool, error) {
		return scope.Update().GetMessage().GetCommand() == command, nil
	}
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
func CallbackDataContains(data string) Condition {
	return func(ctx context.Context, scope *Scope) (bool, error) {
		if scope.Update().IsCallback() {
			cbd := scope.Update().CallbackQuery.Data
			return strings.Contains(cbd, data), nil
		}
		return false, nil
	}
}
