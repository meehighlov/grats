package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func DeleteFriendCallbackQueryHandler(ctx context.Context, event *common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	friendId := params.Id

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, nil)

	if err != nil {
		if _, err := event.ReplyCallbackQuery(ctx, "Что-то пошло не так⚠️ Если проблема повторяется - опишите ее в чате поддержки"); err != nil {
			return err
		}
		event.Logger.Error("error serching friend when deleting: " + err.Error())
		return err
	}

	if len(friends) == 0 {
		event.Logger.Error("not found friend row by id: " + friendId)
		return nil
	}

	friend := friends[0]

	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("⬅️ назад", common.CallInfo(friendId, params.Pagination.Offset, "friend").String()),
		common.NewButton("🗑 удалить", common.CallConfirmDelete(friendId).String()),
	)

	if _, err := event.EditCalbackMessage(ctx, fmt.Sprintf("Удаляем напоминание для %s (%s)?", friend.Name, friend.BirthDay), *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}

func ConfirmDeleteFriendCallbackQueryHandler(ctx context.Context, event *common.Event) error {
	var (
		name     string
		birthDay string
	)

	keyboard := common.NewInlineKeyboard()

	done := false

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		params := common.CallbackFromString(event.GetCallbackQuery().Data)

		friendId := params.Id

		baseFields := db.BaseFields{ID: friendId}
		friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

		if err != nil {
			if _, err := event.ReplyCallbackQuery(ctx, "Что-то пошло не так⚠️ Если проблема повторяется - опишите ее в чате поддержки"); err != nil {
				return err
			}
			event.Logger.Error("error serching friend when deleting: " + err.Error())
			return err
		}

		if len(friends) == 0 {
			return nil
		}

		friend := friends[0]

		err = friend.Delete(ctx, tx)
		if err != nil {
			return err
		}

		done = true

		keyboard.AppendAsStack(common.NewButton("⬅️ к списку др", common.CallList(strconv.Itoa(LIST_START_OFFSET), "<", friend.ChatId, "friend").String()))

		name = friend.Name
		birthDay = friend.BirthDay

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		if _, err := event.EditCalbackMessage(ctx, "Напоминание удалено👋", *keyboard.Murkup()); err != nil {
			return err
		}
		callBackMsg := fmt.Sprintf("Напоминание для %s (%s) удалено🙌", name, birthDay)
		if _, err := event.ReplyCallbackQuery(ctx, callBackMsg); err != nil {
			return err
		}
	}

	return nil
}
