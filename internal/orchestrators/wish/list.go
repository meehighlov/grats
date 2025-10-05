package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) List(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.List(ctx, update)
	})
}

func (o *Orchestrator) WishInfoHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.WishInfoHandler(ctx, update)
	})
}
