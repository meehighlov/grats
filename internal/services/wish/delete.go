package wish

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) DeleteWish(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	wishId := params.ID

	baseFields := entities.BaseFields{ID: wishId}
	wishes, err := s.repositories.Wish.Filter(ctx, &entities.Wish{BaseFields: baseFields})

	if err != nil {
		if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update); err != nil {
			return err
		}
		return err
	}

	if len(wishes) == 0 {
		return nil
	}

	wish := wishes[0]

	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_BACK, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_WISH_INFO, params.Offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_DELETE, s.builders.CallbackDataBuilder.Build(wishId, s.cfg.Constants.CMD_CONFIRM_DELETE_WISH, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
	)

	if _, err := s.clients.Telegram.Edit(ctx, fmt.Sprintf(s.cfg.Constants.DELETE_WISH_CONFIRMATION_TEMPLATE, wish.Name), update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) ConfirmDeleteWish(ctx context.Context, update *telegram.Update) error {
	var (
		wish *entities.Wish
	)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	wishId := params.ID

	baseFields := entities.BaseFields{ID: wishId}
	wishes, err := s.repositories.Wish.Filter(ctx, &entities.Wish{BaseFields: baseFields})

	if err != nil {
		if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update); err != nil {
			return err
		}
		return err
	}

	if len(wishes) == 0 {
		return nil
	}

	wish = wishes[0]

	err = s.repositories.Wish.Delete(ctx, wish)
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

	if _, err := s.clients.Telegram.Edit(ctx, s.cfg.Constants.WISH_DELETED, update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}
