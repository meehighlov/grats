package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) EditPriceHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditPriceHandler(ctx, update)
	})
}

func (o *Orchestrator) SaveEditPriceHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditPriceHandler(ctx, update)
	})
}

func (o *Orchestrator) EditLinkHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditLinkHandler(ctx, update)
	})
}

func (o *Orchestrator) SaveEditLinkHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditLinkHandler(ctx, update)
	})
}

func (o *Orchestrator) EditWishNameHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.EditWishNameHandler(ctx, update)
	})
}

func (o *Orchestrator) SaveEditWishNameHandler(ctx context.Context, update *telegram.Update) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.SaveEditWishNameHandler(ctx, update)
	})
}
