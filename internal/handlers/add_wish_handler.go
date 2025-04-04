package handlers

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	WISH_LIMIT_FOR_USER = 50
	WISH_NAME_MAX_LEN   = 100
)

func AddWishHandler(ctx context.Context, event *common.Event) error {
	wishListId := common.CallbackFromString(event.GetCallbackQuery().Data).Id
	userId := strconv.Itoa(event.GetCallbackQuery().From.Id)

	wishes, err := (&db.Wish{UserId: userId}).Filter(ctx, nil)
	if err != nil {
		event.Logger.Error("error getting wishes: " + err.Error())
		event.Reply(ctx, "Что-то пошло не так⚠️ Если проблема повторяется - опишите ее в чате поддержки")
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

	msg := "✨Введите название желания\n"
	msg += fmt.Sprintf("\n\nМаксимальное количество желаний - %d", WISH_LIMIT_FOR_USER)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}
	event.GetContext().AppendText(wishListId)
	event.GetContext().AppendText(userId)

	event.SetNextHandler("add_save_wish")

	return nil
}

func SaveWish(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()
	userId := strconv.Itoa(message.From.Id)
	wishListId := event.GetContext().GetTexts()[0]

	if wishListId == "" {
		event.Logger.Error(
			"SaveWish",
			"error", "wishListId is empty",
			"userId", userId,
		)
		return nil
	}

	if len(message.Text) > WISH_NAME_MAX_LEN {
		event.Reply(ctx, fmt.Sprintf("Слишком большое имя, максимум - %d символов, попробуйте снова", WISH_NAME_MAX_LEN))
		return nil
	}

	validatedName, err := validateWishName(message.Text)
	if err != nil {
		event.Reply(ctx, "Имя желания содержит недопустимые символы, попробуйте использовать только цифры и буквы")
		return nil
	}

	wish := db.Wish{
		BaseFields: db.NewBaseFields(false),
		Name:       validatedName,
		ChatId:     message.GetChatIdStr(),
		UserId:     userId,
		WishListId: wishListId,
	}

	err = wish.Save(ctx, nil)
	if err != nil {
		return err
	}

	msg := "Желание добавлено 💾"

	if _, err := event.ReplyWithKeyboard(
		ctx,
		msg,
		*buildWishNavigationMarkup(wish.ID, wishListId).Murkup(),
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func buildWishNavigationMarkup(wishId, wishListId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("📂 открыть желание", common.CallWishInfo(wishId, fmt.Sprintf("%d", LIST_START_OFFSET)).String()),
		common.NewButton("➕ новое желание", common.CallAddItem(wishListId, "wish").String()),
		common.NewButton("📋 список желаний", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", wishListId, "wish").String()),
	)

	return keyboard
}

func validateWishName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if len(trimmedName) > WISH_NAME_MAX_LEN {
		return "", fmt.Errorf("название желания не должно превышать %d символов", WISH_NAME_MAX_LEN)
	}

	sanitizedName := html.EscapeString(trimmedName)

	reg := regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}]`)
	sanitizedName = reg.ReplaceAllString(sanitizedName, "")

	return sanitizedName, nil
}
