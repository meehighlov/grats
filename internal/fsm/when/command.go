package when

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type CommandCondition struct {
	command string
}

func Command(command string) *CommandCondition {
	return &CommandCondition{command: command}
}

func (c *CommandCondition) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	return update.GetMessage().GetCommand() == c.command, nil
}
