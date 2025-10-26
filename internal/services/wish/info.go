package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/repositories/models"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
	tgm "github.com/meehighlov/grats/pkg/telegram/models"
)

func (s *Service) WishInfo(ctx context.Context, update *tgm.Update) error {
	var (
		wish *models.Wish
	)
	callbackQuery := update.CallbackQuery

	params := s.builders.CallbackDataBuilder.FromString(callbackQuery.Data)

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, params.ID)
		return err
	})
	if err != nil {
		return err
	}

	offset := params.Offset
	sourceId := wish.WishListId

	viewerId := strconv.Itoa(update.CallbackQuery.From.Id)

	if params.Command == s.cfg.Constants.CMD_SHOW_SWI {
		if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, tgc.WithReplyMurkup(s.buildSharedWishInfoKeyboard(wish, offset, sourceId, viewerId).Murkup())); err != nil {
			return err
		}
	} else {
		if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, tgc.WithReplyMurkup(s.buildWishInfoKeyboard(wish, offset).Murkup())); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildWishInfoKeyboard(wish *models.Wish, offset string) *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsLine(
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_NAME, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_EDIT_WISH_NAME, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_LINK, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_EDIT_LINK, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_PRICE, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_EDIT_PRICE, offset).String()),
	)

	if wish.Link != "" {
		keyboard.AppendAsLine(keyboard.NewURLButton(wish.GetMarketplace(s.GetSiteName), wish.Link))
	}

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_DELETE, s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_DELETE_WISH, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_BACK_TO_WISHLIST, s.builders.CallbackDataBuilder.Build(wish.WishListId, s.cfg.Constants.CMD_LIST, offset).String()),
	)

	return keyboard
}
