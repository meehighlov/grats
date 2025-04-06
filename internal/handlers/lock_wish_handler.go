package handlers

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func ToggleWishLockHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)
	wishId := params.Id
	offset := params.Pagination.Offset
	viewerId := strconv.Itoa(event.GetCallbackQuery().From.Id)

	var wish *db.Wish

	done := false

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		baseFields := db.BaseFields{ID: wishId}
		wishes, err := (&db.Wish{BaseFields: baseFields}).GetWithLock(ctx, tx)
		if err != nil {
			if _, err := event.ReplyCallbackQuery(ctx, common.ERROR_MESSAGE); err != nil {
				return err
			}
			event.Logger.Error("error searching wish when locking: " + err.Error())
			return err
		}

		// wish info was opened too long and expired
		// and owner deleted it
		if len(wishes) == 0 {
			err := refreshWishInfo(
				ctx,
				event,
				wish,
				offset,
				params.Pagination.Direction,
				wish.WishListId,
				viewerId,
			)
			if err != nil {
				return err
			}

			if _, err := event.ReplyCallbackQuery(ctx, "Видимо, желание было удалено🤔 Попробуйте обновить список"); err != nil {
				return err
			}
			return nil
		}

		wish = wishes[0]

		// wish info was opened too long and expired
		// or someone else locked it faster
		if wish.ExecutorId != "" && wish.ExecutorId != viewerId {
			err := refreshWishInfo(
				ctx,
				event,
				wish,
				offset,
				params.Pagination.Direction,
				wish.WishListId,
				viewerId,
			)
			if err != nil {
				return err
			}
			if _, err := event.ReplyCallbackQuery(ctx, "Кто-то уже забронировал это желание, попробуйте обновить список"); err != nil {
				return err
			}
			return nil
		}

		// same user unlocks wish
		if wish.ExecutorId == viewerId {
			viewerId = ""
		}

		wish.ExecutorId = viewerId
		err = wish.Save(ctx, tx)
		if err != nil {
			event.Logger.Error("error saving wish when locking: " + err.Error())
			return err
		}

		done = true

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		err := refreshWishInfo(
			ctx,
			event,
			wish,
			offset,
			params.Pagination.Direction,
			wish.WishListId,
			viewerId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func refreshWishInfo(ctx context.Context, event *common.Event, wish *db.Wish, offset string, direction string, wishListId string, viewerId string) error {
	keyboard := buildSharedWishInfoKeyboard(wish, offset, direction, wishListId, viewerId)
	if _, err := event.EditCalbackMessage(ctx, wish.Info(viewerId), *keyboard.Murkup()); err != nil {
		return err
	}
	return nil
}
