package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) ShareWishListHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Wish.ShareWishListHandler(ctx, update)
}

func (o *Orchestrator) ShowSharedWishlistHandler(ctx context.Context, update *telegram.Update) error {
	// case when called from /start or comes from link
	is_from_start_option := !update.IsCallback()

	if is_from_start_option {
		if err := o.services.User.RegisterOrUpdateUser(ctx, update); err != nil {
			return err
		}
	}

	return o.services.Wish.ShowSharedWishlistHandler(ctx, update)
}
