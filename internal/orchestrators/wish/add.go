package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) AddWish(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.AddWish(ctx, update)
	})
}

func (o *Orchestrator) SaveWish(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveWish(ctx, update)
	})
}
