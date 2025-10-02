package user

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) Start(ctx context.Context, update *telegram.Update) error {
	return o.services.User.Start(ctx, update)
}
