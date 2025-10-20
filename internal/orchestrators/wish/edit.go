package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) EditPrice(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditPrice(ctx, update)
	})
}

func (o *Orchestrator) SaveEditPrice(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditPrice(ctx, update)
	})
}

func (o *Orchestrator) EditLink(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditLink(ctx, update)
	})
}

func (o *Orchestrator) SaveEditLink(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditLink(ctx, update)
	})
}

func (o *Orchestrator) EditWishName(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditWishName(ctx, update)
	})
}

func (o *Orchestrator) SaveEditWishName(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditWishName(ctx, update)
	})
}

func (o *Orchestrator) DeleteLink(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.DeleteLink(ctx, update)
	})
}
