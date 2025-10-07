package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) ToggleWishLockHandler(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	wishId := params.ID
	offset := params.Offset
	viewerId := strconv.Itoa(update.CallbackQuery.From.Id)

	var wish *entities.Wish

	baseFields := entities.BaseFields{ID: wishId}
	wishes, err := s.repositories.Wish.GetWithLock(ctx, &entities.Wish{BaseFields: baseFields})
	if err != nil {
		if _, err := s.clients.Telegram.Reply(ctx, s.constants.ERROR_MESSAGE, update); err != nil {
			return err
		}
		return err
	}

	// wish info was opened too long and expired
	// and owner deleted it
	if len(wishes) == 0 {
		if _, err := s.clients.Telegram.Reply(ctx, s.constants.WISH_REMOVED_TRY_REFRESH, update); err != nil {
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
		if _, err := s.clients.Telegram.Reply(ctx, s.constants.WISH_ALREADY_BOOKED, update); err != nil {
			return err
		}
		return nil
	}

	// same user unlocks wish
	if wish.ExecutorId == viewerId {
		viewerId = ""
	}

	wish.ExecutorId = viewerId
	err = s.repositories.Wish.Save(ctx, wish)
	if err != nil {
		return err
	}

	err = s.refreshWishInfo(
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

	return nil
}

func (s *Service) refreshWishInfo(ctx context.Context, update *telegram.Update, wish *entities.Wish, offset string, wishListId string, viewerId string) error {
	keyboard := s.buildSharedWishInfoKeyboard(wish, offset, wishListId, viewerId)
	if _, err := s.clients.Telegram.Edit(ctx, wish.Info(viewerId), update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}
	return nil
}
