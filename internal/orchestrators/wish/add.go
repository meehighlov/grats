package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) AddWishHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.AddWishHandler(ctx, update)
}

func (o *Orchestrator) SaveWish(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.SaveWish(ctx, update)
}
