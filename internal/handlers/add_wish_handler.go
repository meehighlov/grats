package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

const (
	WISH_LINK_MAX_LEN   = 500
	WISH_LIMIT_FOR_USER = 50
)

func AddWishHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	chatId := common.CallbackFromString(event.GetCallbackQuery().Data).Id
	userId := strconv.Itoa(event.GetCallbackQuery().From.Id)

	wishes, err := (&db.Wish{UserId: userId}).Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("error getting wishes: " + err.Error())
		event.Reply(ctx, "Возникла непредвиденная ошибка, над этим уже работают😔")
		return err
	}

	if len(wishes) >= WISH_LIMIT_FOR_USER {
		event.ReplyCallbackQuery(
			ctx,
			fmt.Sprintf(
				"Достигнут лимит желаний👉👈 Максимальное количество желаний для одного пользователя: %d",
				WISH_LIMIT_FOR_USER,
			),
		)
		return nil
	}

	msg := "✨Введите название желания\n\n"
	msg += fmt.Sprintf("\n\nМаксимальное количество желаний для пользователя: %d", WISH_LIMIT_FOR_USER)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}
	event.GetContext().AppendText(chatId)
	event.GetContext().AppendText(userId)

	event.SetNextHandler("add_save_wish")

	return nil
}

func SaveWish(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	message := event.GetMessage()

	wish := db.Wish{
		BaseFields: db.NewBaseFields(),
		Name:       message.Text,
		ChatId:     message.GetChatIdStr(),
		UserId:     strconv.Itoa(message.From.Id),
		Locked:     "0",
	}

	err := wish.Save(ctx, tx)
	if err != nil {
		return err
	}

	msg := "Желание добавлено 💾"

	if _, err := event.ReplyWithKeyboard(
		ctx,
		msg,
		*buildWishNavigationMarkup(event.GetChatId(), wish.ID).Murkup(),
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func buildWishNavigationMarkup(chatId string, wishId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("➕ новое желание", common.CallAddItem(chatId, "wish").String()),
		common.NewButton("💰 обновить цену", common.CallEditPrice(wishId).String()),
		common.NewButton("🔗 обновить ссылку", common.CallEditLink(wishId).String()),
		common.NewButton("📋 список желаний", common.CallWishList(chatId).String()),
	)

	return keyboard
}
