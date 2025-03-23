package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func FriendInfoCallbackQueryHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching event info: " + err.Error())
		return err
	}

	friend := friends[0]

	// todo take from db
	friendTimezone := "мск"

	emoji, zodiacName := friend.GetZodiacSign()

	msgLines := []string{
		fmt.Sprintf("✨ %s", friend.Name),
		fmt.Sprintf("🗓 %s", friend.BirthDay),
		fmt.Sprintf("%s %s", emoji, zodiacName),
		// todo add info abount bound chat
		fmt.Sprintf("🔔 Напомню %s в полночь (по %s)", *friend.GetNotifyAt(), friendTimezone),
	}

	if friend.IsTodayBirthday() {
		msgLines = append(msgLines, fmt.Sprintf("🥳 Сегодня %s празднует день рождения", friend.Name))
	} else {
		if friend.IsThisMonthAfterToday() {
			msgLines = append(msgLines, fmt.Sprintf("🕒 У %s скоро день рождения", friend.Name))
		}
	}

	msg := strings.Join(msgLines, "\n\n")

	offset := params.Pagination.Offset

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		*common.NewButton("⬅️ к списку др", common.CallList(offset, "<", friend.ChatId).String()),
		*common.NewButton("✏️ редактировать имя", common.CallEditName(params.Id).String()),
		*common.NewButton("📅 редактировать др", common.CallEditBirthday(params.Id).String()),
		*common.NewButton("🗑 удалить", common.CallDelete(params.Id, params.Pagination.Offset).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, msg, *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}
