package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/models"
)

func (s *Service) ToggleWishLock(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	wishId := params.ID
	offset := params.Offset
	viewerId := strconv.Itoa(update.CallbackQuery.From.Id)

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		baseFields := models.BaseFields{ID: wishId}
		wishes, err := s.repositories.Wish.GetWithLock(ctx, &models.Wish{BaseFields: baseFields})
		if err != nil {
			if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update); err != nil {
				return err
			}
			return err
		}

		// wish info was opened too long and expired
		// and owner deleted it
		if len(wishes) == 0 {
			if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_REMOVED_TRY_REFRESH, update); err != nil {
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
				update,
				wish,
				offset,
				wish.WishListId,
				viewerId,
			)
			if err != nil {
				return err
			}
			if _, err := s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_ALREADY_BOOKED, update); err != nil {
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
		update,
		wish,
		offset,
		wish.WishListId,
		viewerId,
	)
}

func (s *Service) refreshWishInfo(ctx context.Context, update *telegram.Update, wish *models.Wish, offset string, wishListId string, viewerId string) error {
	keyboard := s.buildSharedWishInfoKeyboard(wish, offset, wishListId, viewerId)
	if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}
	return nil
}
