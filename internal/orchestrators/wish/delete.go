package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) DeleteWishCallbackQueryHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.DeleteWishCallbackQueryHandler(ctx, update)
}

func (o *Orchestrator) ConfirmDeleteWishCallbackQueryHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.ConfirmDeleteWishCallbackQueryHandler(ctx, update)
}
