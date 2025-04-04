package handlers

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	FRIEND_NAME_MAX_LEN = 50
	EMPTY_CHAT_ID       = "empty"

	FRIEND_LIMIT_FOR_CHAT = 50
)

func AddToChatHandler(ctx context.Context, event *common.Event) error {
	chatId := common.CallbackFromString(event.GetCallbackQuery().Data).Id
	friends, err := (&db.Friend{ChatId: chatId}).Filter(ctx, nil)
	if err != nil {
		event.Logger.Error("error getting friends: " + err.Error())
		return err
	}

	if len(friends) >= FRIEND_LIMIT_FOR_CHAT {
		event.ReplyCallbackQuery(
			ctx,
			fmt.Sprintf(
				"Достигнут лимит напоминаний👉👈 Максимальное количество напоминаний в одном чате: %d",
				FRIEND_LIMIT_FOR_CHAT,
			),
		)
		return nil
	}

	msg := "Введите имя именинника✨\n\nнапример 👉 Райан Гослинг"
	msg += fmt.Sprintf("\n\nМаксимальное количество напоминаний в одном чате: %d", FRIEND_LIMIT_FOR_CHAT)

	if _, err := event.ReplyCallbackQuery(ctx, msg); err != nil {
		return err
	}
	event.GetContext().AppendText(chatId)

	event.SetNextHandler("add_enter_bd")

	return nil
}

func EnterBirthday(ctx context.Context, event *common.Event) error {
	friendName := strings.TrimSpace(event.GetMessage().Text)

	validatedName, err := validateFriendName(friendName)
	if err != nil {
		if _, err := event.Reply(ctx, err.Error()); err != nil {
			return err
		}
		event.SetNextHandler("add_enter_bd")
		return nil
	}

	event.GetContext().AppendText(validatedName)

	msg := "Введите дату рождения✨\n\nформат 👉 день.месяц[.год]\n\nнапример 👉 12.11.1980 или 12.11"

	if _, err := event.Reply(ctx, msg); err != nil {
		return err
	}

	event.SetNextHandler("add_save_friend")

	return nil
}

func SaveFriend(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()
	chatContext := event.GetContext()

	if err := validateBirthdaty(message.Text); err != nil {
		errMsg := "Дата не попадает под формат🤔\n\nВведите дату иначе🙌"
		if _, err := event.Reply(ctx, errMsg); err != nil {
			return err
		}
		event.SetNextHandler("add_save_friend")
		return nil
	}

	chatContext.AppendText(message.Text)
	data := chatContext.GetTexts()
	chatid, name, bd := data[0], data[1], data[2]

	friend := db.Friend{
		BaseFields: db.NewBaseFields(false),
		Name:       name,
		BirthDay:   bd,
		UserId:     strconv.Itoa(message.From.Id),
		ChatId:     chatid,
	}

	friend.RenewNotifayAt()

	err := friend.Save(ctx, nil)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("День рождения для %s добавлен 💾\n\nНапомню о нем %s🔔", name, *friend.GetNotifyAt())

	if strings.Contains(chatid, "-") {
		chatTitle := "чат"
		chatFullInfo, err := event.GetChat(ctx, chatid)
		if err != nil {
			return err
		}
		if chatFullInfo.Id != 0 {
			chatTitle = fmt.Sprintf("чат %s", chatFullInfo.Title)
		}

		msg = fmt.Sprintf("День рождения для %s добавлен в %s 💾\n\nПришлю напоминание в чат %s🔔", name, chatTitle, *friend.GetNotifyAt())
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		msg,
		*buildNavigationMarkup(chatid).Murkup(),
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func validateBirthdaty(birtday string) error {
	month := "01"
	day := "02"
	format_wo_year := fmt.Sprintf("%s.%s", day, month)
	format_with_year := fmt.Sprintf("%s.%s.2006", day, month)

	format := ""

	parts := strings.Split(birtday, ".")
	if len(parts) == 3 {
		format = format_with_year
	} else {
		format = format_wo_year
	}

	_, err := time.Parse(format, birtday)

	if err != nil {
		return err
	}

	return nil
}

func validateFriendName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", fmt.Errorf("имя не может быть пустым")
	}

	if len(trimmedName) > FRIEND_NAME_MAX_LEN {
		return "", fmt.Errorf("имя не должно превышать %d символов", FRIEND_NAME_MAX_LEN)
	}

	sanitizedName := html.EscapeString(trimmedName)

	reg := regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}]`)
	sanitizedName = reg.ReplaceAllString(sanitizedName, "")

	return sanitizedName, nil
}

func buildNavigationMarkup(chatId string) *common.InlineKeyboard {
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsStack(
		common.NewButton("➕ добавить еще", common.CallAddItem(chatId, "friend").String()),
		common.NewButton("📋 список др", common.CallList(fmt.Sprintf("%d", LIST_START_OFFSET), ">", chatId, "friend").String()),
	)

	return keyboard
}
