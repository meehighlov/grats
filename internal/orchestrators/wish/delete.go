package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) DeleteWish(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.DeleteWish(ctx, update)
	})
}

func (o *Orchestrator) ConfirmDeleteWish(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.ConfirmDeleteWish(ctx, update)
	})
}
