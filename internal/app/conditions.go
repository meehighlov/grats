package app

import (
	"context"
	"strings"

	"github.com/meehighlov/grats/pkg/telegram"
)

func SupportReplyCondition(supportChatId string) telegram.Condition {
	return func(ctx context.Context, scope *telegram.Scope) (bool, error) {
		if scope.Update().GetMessage() != nil &&
			scope.Update().GetMessage().GetChatIdStr() == supportChatId &&
			scope.Update().GetMessage().IsReply() {
			return true, nil
		}

		return false, nil
	}
}

func ShowSharedListCondition() func(context.Context, *telegram.Scope) (bool, error) {
	return func(ctx context.Context, scope *telegram.Scope) (bool, error) {
		command := scope.Update().GetMessage().GetCommand()

		if strings.HasPrefix(command, "/start") {
			idForCommand := strings.TrimSpace(strings.TrimPrefix(command, "/start"))

			if idForCommand != "" && strings.HasPrefix(idForCommand, "wl") {
				return true, nil
			}
		}

		return false, nil
	}
}
