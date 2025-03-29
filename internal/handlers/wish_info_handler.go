package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func WishInfoHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	wishes, err := (&db.Wish{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching wish info: " + err.Error())
		return err
	}

	if len(wishes) == 0 {
		if _, err := event.ReplyCallbackQuery(ctx, "Желание не найдено😔"); err != nil {
			return err
		}
		return nil
	}

	wish := wishes[0]

	msgLines := []string{
		"✨ Информация о желании",
	}

	if wish.Link != "" {
		msgLines = append(msgLines, fmt.Sprintf("🔗 Ссылка: %s", wish.Link))
	} else {
		msgLines = append(msgLines, "🔗 Ссылка: отсутствует")
	}

	if wish.OzonLink != "" {
		msgLines = append(msgLines, fmt.Sprintf("🟣 OZON: %s", wish.OzonLink))
	}

	if wish.WbLink != "" {
		msgLines = append(msgLines, fmt.Sprintf("🟡 Wildberries: %s", wish.WbLink))
	}

	if wish.Price != "" {
		msgLines = append(msgLines, fmt.Sprintf("💰 Цена: %s", wish.Price))
	}

	if wish.Locked != "" {
		msgLines = append(msgLines, fmt.Sprintf("🔒 Забронировано: %s", wish.Locked))
	} else {
		msgLines = append(msgLines, "🔓 Статус: не забронировано")
	}

	msg := strings.Join(msgLines, "\n\n")

	offset := params.Pagination.Offset

	keyboard := common.NewInlineKeyboard()

	// todo add link button
	keyboard.AppendAsStack(
		common.NewButton("💰 обновить цену", common.CallEditPrice(wish.ID).String()),
		common.NewButton("🔗 обновить ссылку", common.CallEditLink(wish.ID).String()),
		common.NewButton("⬅️ к списку желаний", common.CallWishList(wish.ChatId).String()),
		common.NewButton("🗑 удалить", common.CallDeleteWish(params.Id, offset).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, msg, *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}
