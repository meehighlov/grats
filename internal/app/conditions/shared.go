package conditions

import (
	"context"
	"strings"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type SharedCondition struct{}

func ShowSharedListCondition() *SharedCondition {
	return &SharedCondition{}
}

func (c *SharedCondition) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	command := update.GetMessage().GetCommand()

	if strings.HasPrefix(command, "/start") {
		idForCommand := strings.TrimSpace(strings.TrimPrefix(command, "/start"))

		if idForCommand != "" && strings.HasPrefix(idForCommand, "wl") {
			return true, nil
		}
	}

	return false, nil
}
