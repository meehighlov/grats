package when

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type CommandCondition struct {
	command string
}

func Command(command string) *CommandCondition {
	return &CommandCondition{command: command}
}

func (c *CommandCondition) Check(ctx context.Context, update *models.Update) (bool, error) {
	return update.GetMessage().GetCommand() == c.command, nil
}
