package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) ToggleWishLockHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.ToggleWishLockHandler(ctx, update)
}
