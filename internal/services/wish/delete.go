package wish

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/pkg/telegram"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) DeleteWish(ctx context.Context, scope *telegram.Scope) error {
	var (
		wish *models.Wish
	)

	params := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)

	wishId := params.ID

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if _, err := scope.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE); err != nil {
			return err
		}
		return err
	}

	keyboard := scope.Keyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_BACK, scope.CallbackData().Build(wish.ID, s.cfg.Constants.CMD_WISH_INFO, params.Offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_DELETE, scope.CallbackData().Build(wishId, s.cfg.Constants.CMD_CONFIRM_DELETE_WISH, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
	)

	if _, err := scope.Edit(ctx, fmt.Sprintf(s.cfg.Constants.DELETE_WISH_CONFIRMATION_TEMPLATE, wish.Name), tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) ConfirmDeleteWish(ctx context.Context, scope *telegram.Scope) error {
	var (
		wish *models.Wish
	)

	params := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)

	wishId := params.ID

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)

		if err != nil {
			if _, err := scope.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE); err != nil {
				return err
			}
			return err
		}

		err = s.repositories.Wish.Delete(ctx, wish)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	keyboard := scope.Keyboard()
	keyboard.AppendAsStack(
		keyboard.NewButton(
			s.cfg.Constants.BTN_BACK_TO_WISHLIST,
			scope.CallbackData().Build(
				wish.WishListId,
				s.cfg.Constants.CMD_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()))

	if _, err := scope.Edit(ctx, s.cfg.Constants.WISH_DELETED, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}
