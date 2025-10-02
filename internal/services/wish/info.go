package wish

import (
	"context"
	"strconv"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) WishInfoHandler(ctx context.Context, update *telegram.Update) error {
	callbackQuery := update.CallbackQuery

	params := s.builders.CallbackDataBuilder.FromString(callbackQuery.Data)

	baseFields := entities.BaseFields{ID: params.ID}
	wishes, err := s.repositories.Wish.Filter(ctx, nil, &entities.Wish{BaseFields: baseFields})
	if err != nil {
		return err
	}

	if len(wishes) == 0 {
		if _, err := s.clients.Telegram.Reply(ctx, s.constants.WISH_WAS_DELETED, update); err != nil {
			return err
		}
		return nil
	}

	wish := wishes[0]

	offset := params.Offset
	sourceId := wish.WishListId

	viewerId := strconv.Itoa(update.CallbackQuery.From.Id)

	if params.Command == s.constants.CMD_SHOW_SWI {
		if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, telegram.WithReplyMurkup(s.buildSharedWishInfoKeyboard(wish, offset, sourceId, viewerId).Murkup())); err != nil {
			return err
		}
	} else {
		if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, telegram.WithReplyMurkup(s.buildWishInfoKeyboard(wish, offset).Murkup())); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildWishInfoKeyboard(wish *entities.Wish, offset string) *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BTN_EDIT_NAME, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_EDIT_WISH_NAME, offset).String()),
		keyboard.NewButton(s.constants.BTN_EDIT_LINK, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_EDIT_LINK, offset).String()),
		keyboard.NewButton(s.constants.BTN_EDIT_PRICE, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_EDIT_PRICE, offset).String()),
	)

	if wish.Link != "" {
		keyboard.AppendAsLine(keyboard.NewURLButton(wish.GetMarketplace(), wish.Link))
	}

	keyboard.AppendAsStack(
		keyboard.NewButton(s.constants.BTN_DELETE, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_DELETE_WISH, offset).String()),
		keyboard.NewButton(s.constants.BTN_BACK_TO_WISHLIST, s.builders.CallbackDataBuilder.Build(wish.WishListId, s.constants.CMD_LIST, offset).String()),
	)

	return keyboard
}
