package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

	msg := "Введите ссылку на желание✨\n\nнапример 👉 https://example.com/wish"
	msg += fmt.Sprintf("\n\nМаксимальное количество желаний для пользователя: %d", WISH_LIMIT_FOR_USER)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}
	event.GetContext().AppendText(chatId)
	event.GetContext().AppendText(userId)

	event.SetNextHandler("add_enter_ozon_link")

	return nil
}

func EnterOzonLink(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	wishLink := strings.TrimSpace(event.GetMessage().Text)

	if len(wishLink) > WISH_LINK_MAX_LEN {
		if _, err := event.Reply(ctx, fmt.Sprintf("Ссылка не должна превышать %d символов", WISH_LINK_MAX_LEN)); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_ozon_link")
		return nil
	}

	event.GetContext().AppendText(wishLink)

	msg := "Введите ссылку на OZON (если есть)✨\n\nнапример 👉 https://ozon.ru/product/123 или введите '-' если её нет"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_enter_wb_link")

	return nil
}

func EnterWbLink(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	ozonLink := strings.TrimSpace(event.GetMessage().Text)

	if ozonLink == "-" {
		ozonLink = ""
	}

	if len(ozonLink) > WISH_LINK_MAX_LEN {
		if _, err := event.Reply(ctx, fmt.Sprintf("Ссылка не должна превышать %d символов", WISH_LINK_MAX_LEN)); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_wb_link")
		return nil
	}

	event.GetContext().AppendText(ozonLink)

	msg := "Введите ссылку на Wildberries (если есть)✨\n\nнапример 👉 https://wildberries.ru/catalog/123 или введите '-' если её нет"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_enter_price")

	return nil
}

func EnterPrice(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	wbLink := strings.TrimSpace(event.GetMessage().Text)

	if wbLink == "-" {
		wbLink = ""
	}

	if len(wbLink) > WISH_LINK_MAX_LEN {
		if _, err := event.Reply(ctx, fmt.Sprintf("Ссылка не должна превышать %d символов", WISH_LINK_MAX_LEN)); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_price")
		return nil
	}

	event.GetContext().AppendText(wbLink)

	msg := "Укажите примерную цену (если известна)✨\n\nнапример 👉 1000₽ или введите '-' если цена неизвестна"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_save_wish")

	return nil
}

func SaveWish(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	message := event.GetMessage()
	chatContext := event.GetContext()

	price := strings.TrimSpace(message.Text)
	if price == "-" {
		price = ""
	}

	chatContext.AppendText(price)
	data := chatContext.GetTexts()
	chatId, userId, link, ozonLink, wbLink, priceVal := data[0], data[1], data[2], data[3], data[4], data[5]

	wish := db.Wish{
		BaseFields: db.NewBaseFields(),
		ChatId:     chatId,
		UserId:     userId,
		Link:       link,
		OzonLink:   ozonLink,
		WbLink:     wbLink,
		Locked:     "",
		Price:      priceVal,
	}

	err := wish.Save(ctx, tx)
	if err != nil {
		return err
	}

	msg := "Желание успешно добавлено 💾"

	if strings.Contains(chatId, "-") {
		chatTitle := "чат"
		chatFullInfo, err := event.GetChat(ctx, chatId)
		if err != nil {
			return err
		}
		if chatFullInfo.Id != 0 {
			chatTitle = fmt.Sprintf("чат %s", chatFullInfo.Title)
		}

		msg = fmt.Sprintf("Желание добавлено в %s 💾", chatTitle)
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		msg,
		*buildWishNavigationMarkup(chatId).Murkup(),
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func buildWishNavigationMarkup(chatId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("➕ добавить еще", common.CallAddItem(chatId, "wish").String()),
		common.NewButton("📋 список желаний", common.CallWishList(chatId).String()),
	)

	return keyboard
}
