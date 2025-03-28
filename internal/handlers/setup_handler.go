package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func SetupHandler(ctx context.Context, event *common.Event, _ *gorm.DB) error {
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

	listButton := common.NewButton("🎂 Личные напоминания", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", chatId).String())
	groupButton := common.NewButton("👥 Групповые чаты", common.CallChatList().String())
	supportButton := common.NewButton("💬 Чат с поддержкой", common.CallSupport(chatId).String())

	// Добавляем кнопку для добавления бота в чат
	cfg := config.Cfg()
	addBotButton := common.NewAddBotToChatURLButton("➕ Добавить бота в чат", cfg.BotName)

	keyboard.AppendAsStack(*listButton, *groupButton, *supportButton)
	keyboard.AppendAsLine(*addBotButton)

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

func SetupFromGroupHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	chatId := event.GetMessage().GetChatIdStr()

	friends, err := (&db.Friend{ChatId: chatId}).Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("error fetching friends", "error", err.Error())
		return err
	}

	if len(friends) == 0 {
		event.Reply(ctx, "Нет информации о ближайших др🙌")
		return nil
	}

	nearest := []*db.Friend{}
	for _, friend := range friends {
		if friend.IsThisMonthAfterToday() || friend.IsTodayBirthday() {
			nearest = append(nearest, friend)
		}
	}

	if len(nearest) == 0 {
		event.Reply(ctx, "Нет др в этом месяце✨")
		return nil
	}

	msg := ""
	for _, friend := range nearest {
		if friend.IsTodayBirthday() {
			msg += fmt.Sprintf("🥳 др сегодня  %s - %s", friend.Name, friend.BirthDay)
		} else {
			msg += fmt.Sprintf("🕒 др в этом месяце %s - %s", friend.Name, friend.BirthDay)
		}
		msg += "\n"
	}

	event.Reply(ctx, msg)

	return nil
}
