package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func CommandListHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()

	// at some point it is possible to use /command in group chat
	// so block this action
	if strings.HasSuffix(message.Chat.Type, "group") {
		return nil
	}

	keyboard := common.NewInlineKeyboard()

	chatId := event.GetMessage().GetChatIdStr()
	if event.GetCallbackQuery().Id != "" {
		chatId = strconv.Itoa(event.GetCallbackQuery().From.Id)
	}

	// while we have only one wishlist per chat
	// we can just get the first one by chat id
	wishLists, err := (&db.WishList{ChatId: chatId}).Filter(ctx, nil)
	if err != nil {
		event.Logger.Error("CommandListHandler error getting wishlists", "chatId", chatId, "error", err.Error())
		return err
	}

	if len(wishLists) == 0 {
		event.Logger.Error("CommandListHandler", "chatId", chatId, "error", "no wishlists found")
		return nil
	}

	wishListId := wishLists[0].ID

	listButton := common.NewButton("🎂 Личные напоминания", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", chatId, "friend").String())
	groupButton := common.NewButton("👥 Групповые чаты", common.CallChatList().String())
	wishButton := common.NewButton("🎁 Список желаний", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", wishListId, "wish").String())
	supportButton := common.NewButton("💬 Чат с поддержкой", common.CallSupport(chatId).String())

	keyboard.AppendAsStack(listButton, groupButton, wishButton, supportButton)

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			"Это список моих комманд🙌",
			*keyboard.Murkup(),
		); err != nil {
			return err
		}
		return nil
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		"Это список моих комманд🙌",
		*keyboard.Murkup(),
	); err != nil {
		return err
	}

	return nil
}
