package conditions

import (
	"context"
	"strings"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type SharedCondition struct{}

func ShowSharedListCondition() *SharedCondition {
	return &SharedCondition{}
}

func (c *SharedCondition) Check(ctx context.Context, update *models.Update) (bool, error) {
	command := update.GetMessage().GetCommand()

	if strings.HasPrefix(command, "/start") {
		idForCommand := strings.TrimSpace(strings.TrimPrefix(command, "/start"))

		if idForCommand != "" && strings.HasPrefix(idForCommand, "wl") {
			return true, nil
		}
	}

	return false, nil
}
