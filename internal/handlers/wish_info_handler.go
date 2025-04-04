package handlers

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func WishInfoHandler(ctx context.Context, event *common.Event) error {
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, nil)

	if err != nil {
		event.Logger.Error("error during fetching wish info: " + err.Error())
		return err
	}

	if len(wishes) == 0 {
		if _, err := event.ReplyCallbackQuery(ctx, "Желание было удалено, попробуйте обновить список"); err != nil {
			return err
		}
		return nil
	}

	wish := wishes[0]

	offset := params.Pagination.Offset
	direction := params.Pagination.Direction
	sourceId := wish.WishListId

	viewerId := strconv.Itoa(event.GetCallbackQuery().From.Id)

	if params.Command == "show_swi" {
		if _, err := event.EditCalbackMessage(ctx, wish.Info(viewerId), *buildSharedWishInfoKeyboard(wish, offset, direction, sourceId, viewerId).Murkup()); err != nil {
			return err
		}
	} else {
		if _, err := event.EditCalbackMessage(ctx, wish.Info(viewerId), *buildWishInfoKeyboard(wish, offset, direction, sourceId).Murkup()); err != nil {
			return err
		}
	}

	return nil
}

func buildWishInfoKeyboard(wish *db.Wish, offset, direction, sourceId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsLine(
		common.NewButton("✏️ название", common.CallEditWishName(wish.ID).String()),
		common.NewButton("✏️ ссылка", common.CallEditLink(wish.ID).String()),
		common.NewButton("✏️ цена", common.CallEditPrice(wish.ID).String()),
	)

	if wish.Link != "" {
		keyboard.AppendAsLine(common.NewURLButton(wish.GetMarketplace(), wish.Link))
	}

	keyboard.AppendAsStack(
		common.NewButton("🗑️ удалить", common.CallDeleteWish(wish.ID, offset).String()),
		common.NewButton("⬅️ к списку желаний", common.CallList(offset, direction, sourceId, "wish").String()),
	)

	return keyboard
}

func buildSharedWishInfoKeyboard(
	wish *db.Wish,
	offset,
	direction,
	sourceId string,
	viewerId string,
) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	if wish.Link != "" {
		keyboard.AppendAsLine(common.NewURLButton(wish.GetMarketplace(), wish.Link))
	}

	if wish.ExecutorId != "" {
		if wish.ExecutorId == viewerId {
			keyboard.AppendAsLine(common.NewButton("✖️ отменить бронь", common.CallToggleWishLock(wish.ID, offset).String()))
		}
		// has executor but it's not viewer - not show lock button
	} else {
		keyboard.AppendAsLine(common.NewButton("🎁 забронировать", common.CallToggleWishLock(wish.ID, offset).String()))
	}

	keyboard.AppendAsStack(
		common.NewButton("⬅️ к списку желаний", common.CallSharedWishList(offset, direction, sourceId, "wish_list").String()),
	)

	return keyboard
}
