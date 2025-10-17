package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"gorm.io/gorm"
)

func (o *Orchestrator) ShareWishList(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)
		return o.services.Wish.ShareWishList(ctx, update)
	})
}

func (o *Orchestrator) ShowSharedWishlist(ctx context.Context, update *telegram.Update) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, o.cfg.TxKey, tx)

		// case when called from /start or comes from link
		isFromStartOption := !update.IsCallback()

		if isFromStartOption {
			if err := o.services.User.RegisterOrUpdateUser(ctx, update); err != nil {
				return err
			}
		}

		return o.services.Wish.ShowSharedWishlist(ctx, update)
	})
}
