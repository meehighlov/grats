package wish

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/repositories/models"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
	tgm "github.com/meehighlov/grats/pkg/telegram/models"
)

func (s *Service) DeleteWish(ctx context.Context, update *tgm.Update) error {
	var (
		wish *models.Wish
	)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	wishId := params.ID

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update); err != nil {
			return err
		}
		return err
	}

	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_BACK, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_WISH_INFO, params.Offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_DELETE, s.builders.CallbackDataBuilder.Build(wishId, s.cfg.Constants.CMD_CONFIRM_DELETE_WISH, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
	)

	if _, err := s.clients.Telegram.Edit(ctx, fmt.Sprintf(s.cfg.Constants.DELETE_WISH_CONFIRMATION_TEMPLATE, wish.Name), update, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) ConfirmDeleteWish(ctx context.Context, update *tgm.Update) error {
	var (
		wish *models.Wish
	)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	wishId := params.ID

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)

		if err != nil {
			if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update); err != nil {
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

	keyboard := s.builders.KeyboardBuilder.NewKeyboard()
	keyboard.AppendAsStack(
		keyboard.NewButton(
			s.cfg.Constants.BTN_BACK_TO_WISHLIST,
			s.builders.CallbackDataBuilder.Build(
				wish.WishListId,
				s.cfg.Constants.CMD_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()))

	if _, err := s.clients.Telegram.Edit(ctx, s.cfg.Constants.WISH_DELETED, update, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}
