package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

const (
	SHARE_WISHLIST_LINK_TEMPLATE = "https://t.me/%s?start=wl%s"
)

func ShareWishListHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	wishListId := params.Id

	wishlist, err := (&db.WishList{BaseFields: db.BaseFields{ID: wishListId}}).Filter(ctx, nil)
	if err != nil {
		return err
	}

	if len(wishlist) == 0 {
		return fmt.Errorf("wishlist not found")
	}

	botName := config.Cfg().BotName
	shareLink := fmt.Sprintf(SHARE_WISHLIST_LINK_TEMPLATE, botName, wishlist[0].ID)

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsLine(
		common.NewShareLinkButton("📤 поделиться", shareLink, "Мой wishlist✨"),
		common.NewCopyButton("🔗 ссылка", shareLink),
	)
	keyboard.AppendAsLine(
		common.NewButton("⬅️ к списку желаний", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", wishlist[0].ID, "wish").String()),
	)

	shareMessage := ("Поделитесь своим вишлистом! Вы так же можете разместить ссылку на вишлист в соцсетях" +
		"\n\n- При переходе по ссылке откроется чат с grats и вишлист будет прислан в виде нового сообщения\n" +
		"- Если пользователь ранее не использовал grats, то ему потребуется лишь нажать start")

	if _, err := event.EditCalbackMessage(
		ctx,
		shareMessage,
		*keyboard.Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func ShowSharedWishlistHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()

	calledFromCallback := false
	wishlistId := strings.TrimPrefix(message.Text, "/start wl")
	offset := LIST_START_OFFSET
	direction := "<"
	if event.GetCallbackQuery().Id != "" {
		calledFromCallback = true
		params := common.CallbackFromString(event.GetCallbackQuery().Data)
		wishlistId = params.Id
		offset, _ = strconv.Atoi(params.Pagination.Offset)
		if offset == 0 {
			offset = LIST_START_OFFSET
		}
		direction = params.Pagination.Direction
	}

	// case when called from /start or comes from link
	if !calledFromCallback {
		err := RegisterOrUpdateUser(ctx, event)
		if err != nil {
			event.Logger.Error("start error registering user", "chatId", message.GetChatIdStr(), "error", err.Error())
			return err
		}
	}

	wishes, err := (&db.Wish{WishListId: wishlistId}).Filter(ctx, nil)
	if err != nil {
		event.Logger.Error(
			"StartHandler - Shared wishlist",
			"error", "error getting wishes",
			"details", err.Error(),
			"wl_id", wishlistId,
			"chatId", message.GetChatIdStr(),
		)
		event.Reply(ctx, "Не удалось загрузить желания пользователя😔")
		return err
	}

	if len(wishes) == 0 {
		// TODO: update message when user list is empty
		return nil
	}

	// user used his own link
	// just send greeting
	if wishes[0].UserId == strconv.Itoa(event.GetMessage().From.Id) && config.Cfg().IsProd() {
		event.Reply(ctx, "Снова привет👋")
		return nil
	}

	userInfo, _ := event.GetChatMember(ctx, wishes[0].UserId)
	name := userInfo.Result.User.Username
	if name == "" {
		name = userInfo.Result.User.FirstName
	}
	header := fmt.Sprintf("✨Вишлист %s", "@"+name)

	var entities []common.PaginatedEntity
	for _, wish := range wishes {
		entities = append(entities, wish)
	}

	offset = common.GetOffsetByDirection(direction, offset, entities, LIST_PAGINATION_SHIFT)

	keyboard := buildSharedWishlistMarkup(entities, LIST_LIMIT, offset, wishlistId, "wish")

	if calledFromCallback {
		if _, err := event.EditCalbackMessage(
			ctx,
			header,
			*keyboard.Murkup(),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		header,
		*keyboard.Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func buildSharedWishlistMarkup(entities []common.PaginatedEntity, limit, offset int, sourceId, table string) *common.InlineKeyboard {
	entityListAsButtons := common.BuildEntityButtons(entities, limit, offset, func(id string, offset int) string {
		return common.CallSharedWishInfo(id, strconv.Itoa(offset)).String()
	})

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsLine(
		common.NewButton("🔄", common.CallSharedWishList(strconv.Itoa(offset), "<", sourceId, "wish_list").String()),
	)
	keyboard.AppendAsStack(entityListAsButtons...)

	common.AppendControlButtons(keyboard, len(entities), limit, offset, sourceId, table, common.CallSharedWishList, LIST_START_OFFSET)
	return keyboard
}
