package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) AddWishHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.AddWishHandler(ctx, update)
	})
}

func (o *Orchestrator) SaveWish(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveWish(ctx, update)
	})
}
