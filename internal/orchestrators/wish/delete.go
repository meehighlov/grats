package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) DeleteWishCallbackQueryHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.DeleteWishCallbackQueryHandler(ctx, update)
	})
}

func (o *Orchestrator) ConfirmDeleteWishCallbackQueryHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.ConfirmDeleteWishCallbackQueryHandler(ctx, update)
	})
}
