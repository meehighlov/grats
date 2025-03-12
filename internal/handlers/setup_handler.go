package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func SetupHandler(ctx context.Context, event *common.Event, _ *sql.Tx) error {
	keyboard := common.NewInlineKeyboard()

	chatId := event.GetMessage().GetChatIdStr()
	if event.GetCallbackQuery().Id != "" {
		chatId = strconv.Itoa(event.GetCallbackQuery().From.Id)
	}

	listButton := common.NewButton("🎂 Личные напоминания", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", chatId).String())
	groupButton := common.NewButton("👥 Групповые чаты", common.CallChatList().String())

	keyboard.AppendAsStack(*listButton, *groupButton)

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

func SetupFromGroupHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	tgChatId := event.GetMessage().GetChatIdStr()

	chat := db.Chat{
		TGChatId: tgChatId,
	}
	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		event.Logger.Error("error fetching chats", "error", err.Error())
		return err
	}

	if len(chats) == 0 {
		event.Reply(ctx, "Нет информации о ближайших др🙌")
		return nil
	}

	friends, err := (&db.Friend{ChatId: chats[0].ID}).Filter(ctx, tx)
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
			nearest = append(nearest, &friend)
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
