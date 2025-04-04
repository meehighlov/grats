package handlers

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func DeleteWishCallbackQueryHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	wishId := params.Id

	baseFields := db.BaseFields{ID: wishId}
	wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, nil)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Что-то пошло не так⚠️ Если проблема повторяется - опишите ее в чате поддержки"); err != nil {
			return err
		}
		event.Logger.Error("error searching wish when deleting: " + err.Error())
		return err
	}

	if len(wishes) == 0 {
		return nil
	}

	wish := wishes[0]

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("⬅️ назад", common.CallWishInfo(wish.ID, params.Pagination.Offset).String()),
		common.NewButton("🗑️ удалить", common.CallConfirmDeleteWish(wishId).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Удалить желание %s?", wish.Name), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ConfirmDeleteWishCallbackQueryHandler(ctx context.Context, event *common.Event) error {
	var (
		wish *db.Wish
	)

	keyboard := common.NewInlineKeyboard()

	done := false

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		params := common.CallbackFromString(event.GetCallbackQuery().Data)

		wishId := params.Id

		baseFields := db.BaseFields{ID: wishId}
		wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, tx)

		if err != nil {
			if _, err := event.ReplyCallbackQuery(ctx, "Что-то пошло не так⚠️ Если проблема повторяется - опишите ее в чате поддержки"); err != nil {
				return err
			}
			event.Logger.Error("error searching wish when deleting: " + err.Error())
			return err
		}

		if len(wishes) == 0 {
			return nil
		}

		wish = wishes[0]

		err = wish.Delete(ctx, tx)
		if err != nil {
			return err
		}

		done = true

		keyboard.AppendAsStack(common.NewButton("⬅️ к списку желаний", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", wish.WishListId, "wish").String()))

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		msg := fmt.Sprintf("Желание %s удалено", wish.Name)
		if _, err := event.EditCalbackMessage(ctx, msg, *keyboard.Murkup()); err != nil {
			return err
		}
	}

	return nil
}
