package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/pkg/telegram"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) ToggleWishLock(ctx context.Context, scope *telegram.Scope) error {
	var (
		wish *models.Wish
	)
	params := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)
	wishId := params.ID
	offset := params.Offset
	viewerId := strconv.Itoa(scope.Update().CallbackQuery.From.Id)

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		baseFields := models.BaseFields{ID: wishId}
		wishes, err := s.repositories.Wish.GetWithLock(ctx, &models.Wish{BaseFields: baseFields})
		if err != nil {
			if _, err := scope.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE); err != nil {
				return err
			}
			return err
		}

		// wish info was opened too long and expired
		// and owner deleted it
		if len(wishes) == 0 {
			if _, err := scope.Reply(ctx, s.cfg.Constants.WISH_REMOVED_TRY_REFRESH); err != nil {
				return err
			}
			return nil
		}

		wish = wishes[0]

		// wish info was opened too long and expired
		// or someone else locked it faster
		if wish.ExecutorId != "" && wish.ExecutorId != viewerId {
			err := s.refreshWishInfo(
				ctx,
				scope,
				wish,
				offset,
				wish.WishListId,
				viewerId,
			)
			if err != nil {
				return err
			}
			if _, err := scope.Reply(ctx, s.cfg.Constants.WISH_ALREADY_BOOKED); err != nil {
				return err
			}
			return nil
		}

		// same user unlocks wish
		if wish.ExecutorId == viewerId {
			viewerId = ""
		}

		wish.ExecutorId = viewerId
		return s.repositories.Wish.Save(ctx, wish)
	})
	if err != nil {
		return err
	}

	return s.refreshWishInfo(
		ctx,
		scope,
		wish,
		offset,
		wish.WishListId,
		viewerId,
	)
}

func (s *Service) refreshWishInfo(ctx context.Context, scope *telegram.Scope, wish *models.Wish, offset string, wishListId string, viewerId string) error {
	keyboard := s.buildSharedWishInfoKeyboard(scope, wish, offset, wishListId, viewerId)
	if _, err := scope.Edit(ctx, wish.Info(viewerId), tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}
	return nil
}
