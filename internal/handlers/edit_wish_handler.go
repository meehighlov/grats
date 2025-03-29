package handlers

import (
	"context"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func EditPriceHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	callbackData := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(ctx, "Введите цену")

	event.GetContext().AppendText(callbackData.Id)
	event.SetNextHandler("edit_price_save")

	return nil
}

func SaveEditPriceHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	message := event.GetMessage()

	wish, err := (&db.Wish{BaseFields: db.BaseFields{ID: event.GetContext().GetTexts()[0]}}).Filter(ctx, tx)
	if err != nil {
		event.Reply(ctx, "Не смогли найти желание")
		return err
	}

	// todo validate price
	wish[0].Price = message.Text

	wish[0].Save(ctx, tx)

	// todo add info button
	event.Reply(ctx, "Цена установлена 💾")

	return nil
}

func EditLinkHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	callbackData := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(ctx, "Введите ссылку")

	event.GetContext().AppendText(callbackData.Id)
	event.SetNextHandler("edit_link_save")

	return nil
}

func SaveEditLinkHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	message := event.GetMessage()

	wish, err := (&db.Wish{BaseFields: db.BaseFields{ID: event.GetContext().GetTexts()[0]}}).Filter(ctx, tx)
	if err != nil {
		event.Reply(ctx, "Не смогли найти желание")
		return err
	}

	// todo validate link
	wish[0].Link = message.Text

	wish[0].Save(ctx, tx)

	// todo add info button
	event.Reply(ctx, "Ссылка установлена 💾")

	return nil
}
