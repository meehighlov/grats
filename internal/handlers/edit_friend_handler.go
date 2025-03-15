package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

func EditNameHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	callbackQuery := event.GetCallbackQuery()
	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching friend info: " + err.Error())
		return err
	}

	friend := friends[0]

	msg := fmt.Sprintf("Введите новое имя для %s", friend.Name)

	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		*common.NewButton("⬅️ назад", common.CallInfo(params.Id, "0").String()),
	)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}

	event.GetContext().AppendText(params.Id)
	event.SetNextHandler("save_edit_name")

	return nil
}

func SaveEditNameHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()
	chatContext := event.GetContext()

	newName := strings.TrimSpace(message.Text)

	if len(newName) > FRIEND_NAME_MAX_LEN {
		if _, err := event.Reply(ctx, fmt.Sprintf("Имя не должно превышать %d символов", FRIEND_NAME_MAX_LEN)); err != nil {
			return err
		}
		event.SetNextHandler("save_edit_name")
		return nil
	}

	friendId := chatContext.GetTexts()[0]

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching friend info: " + err.Error())
		return err
	}

	friend := friends[0]
	friend.Name = newName

	err = friend.Save(ctx, tx)
	if err != nil {
		event.Logger.Error("error saving friend with new name: " + err.Error())
		return err
	}

	msg := "Имя изменено 💾"

	replyWithInfo(ctx, event, friend, msg)

	event.SetNextHandler("")

	return nil
}

func EditBirthdayHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	callbackQuery := event.GetCallbackQuery()
	params := common.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching friend info: " + err.Error())
		return err
	}

	friend := friends[0]

	msg := fmt.Sprintf("Введите новую дату рождения для %s\n\nТекущая дата: %s\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 12.11.1980 или 12.11", friend.Name, friend.BirthDay)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}

	event.GetContext().AppendText(params.Id)
	event.SetNextHandler("save_edit_birthday")

	return nil
}

func SaveEditBirthdayHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	message := event.GetMessage()
	chatContext := event.GetContext()

	newBirthday := strings.TrimSpace(message.Text)

	if err := validateBirthdaty(newBirthday); err != nil {
		errMsg := "Дата не попадает под формат🤔\n\nВведите дату иначе🙌"
		if _, err := event.Reply(ctx, errMsg); err != nil {
			return err
		}
		event.SetNextHandler("save_edit_birthday")
		return nil
	}

	friendId := chatContext.GetTexts()[0]

	baseFields := db.BaseFields{ID: friendId}
	friends, err := (&db.Friend{BaseFields: baseFields}).Filter(ctx, tx)

	if err != nil {
		event.Logger.Error("error during fetching friend info: " + err.Error())
		return err
	}

	friend := friends[0]
	oldBirthday := friend.BirthDay

	msgTemplate := "Дата рождения %s изменена 💾\n\nНапомню о нем %s🔔"

	if strings.EqualFold(newBirthday, oldBirthday) {
		replyWithInfo(ctx, event, friend, fmt.Sprintf(msgTemplate, friend.Name, *friend.GetNotifyAt()))
		event.SetNextHandler("")
		return nil
	}

	friend.BirthDay = newBirthday
	friend.RenewNotifayAt()

	err = friend.Save(ctx, tx)
	if err != nil {
		event.Logger.Error("SaveEditBirthdayHandler", "birthday update error", err.Error())
		return err
	}

	if err := replyWithInfo(ctx, event, friend, fmt.Sprintf(msgTemplate, friend.Name, *friend.GetNotifyAt())); err != nil {
		event.Logger.Error("SaveEditBirthdayHandler", "reply error", err.Error())
		return err
	}

	event.SetNextHandler("")

	return nil
}

func replyWithInfo(
	ctx context.Context,
	event *common.Event,
	friend *db.Friend,
	msg string,
) error {
	keyboard := common.NewInlineKeyboard()
	keyboard.AppendAsStack(
		*common.NewButton(
			fmt.Sprintf("%s %s", friend.Name, friend.BirthDay),
			common.CallInfo(friend.ID, fmt.Sprintf("%d", LIST_START_OFFSET)).String()),
	)

	if _, err := event.ReplyWithKeyboard(ctx, msg, *keyboard.Murkup()); err != nil {
		return err
	}

	return nil
}
