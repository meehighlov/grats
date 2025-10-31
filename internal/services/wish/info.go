package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/pkg/telegram"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) WishInfo(ctx context.Context, scope *telegram.Scope) error {
	var (
		wish *models.Wish
	)
	callbackQuery := scope.Update().CallbackQuery

	params := scope.CallbackData().FromString(callbackQuery.Data)

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, params.ID)
		return err
	})
	if err != nil {
		return err
	}

	offset := params.Offset
	sourceId := wish.WishListId

	viewerId := strconv.Itoa(scope.Update().CallbackQuery.From.Id)

	if params.Command == s.cfg.Constants.CMD_SHOW_SWI {
		if _, err := scope.Edit(ctx, wish.Info(viewerId), tgc.WithReplyMurkup(s.buildSharedWishInfoKeyboard(scope, wish, offset, sourceId, viewerId).Murkup())); err != nil {
			return err
		}
	} else {
		if _, err := scope.Edit(ctx, wish.Info(viewerId), tgc.WithReplyMurkup(s.buildWishInfoKeyboard(scope, wish, offset).Murkup())); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildWishInfoKeyboard(scope *telegram.Scope, wish *models.Wish, offset string) *inlinekeyboard.Builder {
	keyboard := scope.Keyboard()
	callbackData := scope.CallbackData()

	keyboard.AppendAsLine(
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_NAME, callbackData.Build(wish.ID, s.cfg.Constants.CMD_EDIT_WISH_NAME, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_LINK, callbackData.Build(wish.ID, s.cfg.Constants.CMD_EDIT_LINK, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_EDIT_PRICE, callbackData.Build(wish.ID, s.cfg.Constants.CMD_EDIT_PRICE, offset).String()),
	)

	if wish.Link != "" {
		keyboard.AppendAsLine(keyboard.NewURLButton(wish.GetMarketplace(s.GetSiteName), wish.Link))
	}

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_DELETE, callbackData.Build(wish.ID, s.cfg.Constants.CMD_DELETE_WISH, offset).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_BACK_TO_WISHLIST, callbackData.Build(wish.WishListId, s.cfg.Constants.CMD_LIST, offset).String()),
	)

	return keyboard
}
