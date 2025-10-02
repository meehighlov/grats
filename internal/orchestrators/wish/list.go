package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) List(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.List(ctx, update)
}

func (o *Orchestrator) WishInfoHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.WishInfoHandler(ctx, update)
}
