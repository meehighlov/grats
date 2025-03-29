package handlers

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func DeleteWishCallbackQueryHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	wishId := params.Id

	baseFields := db.BaseFields{ID: wishId}
	wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error searching wish when deleting: " + err.Error())
		return err
	}

	if len(wishes) == 0 {
		event.Logger.Error("not found wish row by id: " + wishId)
		return err
	}

	wish := wishes[0]

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("⬅️ назад", common.CallWishList(wish.ChatId).String()),
		common.NewButton("🗑 удалить", common.CallConfirmDeleteWish(wishId).String()),
	)

	linkDisplay := wish.Link
	if linkDisplay == "" {
		linkDisplay = "Желание без ссылки"
	}

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Удаляем желание %s?", linkDisplay), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ConfirmDeleteWishCallbackQueryHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	wishId := params.Id

	baseFields := db.BaseFields{ID: wishId}
	wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error searching wish when deleting: " + err.Error())
		return err
	}

	if len(wishes) == 0 {
		event.Logger.Error("not found wish row by id: " + wishId)
		return err
	}

	wish := wishes[0]

	err = wish.Delete(ctx, tx)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔"); err != nil {
			return err
		}
		event.Logger.Error("error deleting wish: " + err.Error())
		return err
	}

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(common.NewButton("⬅️ к списку желаний", common.CallWishList(wish.ChatId).String()))

	if _, err := event.EditCalbackMessage(ctx, "Желание удалено👋", *keyboard.Murkup()); err != nil {
		return err
	}

	linkDisplay := wish.Link
	if linkDisplay == "" {
		linkDisplay = "Желание без ссылки"
	}
	callBackMsg := fmt.Sprintf("Желание %s удалено🙌", linkDisplay)
	if _, err := event.ReplyCallbackQuery(ctx, callBackMsg); err != nil {
		return err
	}

	return nil
}
