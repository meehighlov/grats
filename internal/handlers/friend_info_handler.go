package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func FriendInfoCallbackQueryHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
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

	markup := [][]map[string]string{
		{
			{
				"text":          "👈 к списку др",
				"callback_data": common.CallList(offset, "<", params.BoundChat).String(),
			},
		},
		{
			{
				"text":          "удалить 👋",
				"callback_data": common.CallDelete(params.Id, params.BoundChat).String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
