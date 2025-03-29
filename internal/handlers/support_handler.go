package handlers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/telegram"
	"gorm.io/gorm"
)

const (
	SUPPORT_WELCOME_MESSAGE = `
	Напишите нам, если у Вас возникли вопросы или предложения%sТакже Вы можете оставить обратную связь о работе grats👋`

	SUPPORT_INFO_MESSAGE = `
	Пришлите Ваш текст обращения%sПожалуйста, будьте вежливы в чате, мы стараемся для Вас😌
	`

	SUPPORT_THANKS_MESSAGE = `
	Спасибо за обращение! Как только поддержка обработает Ваш запрос, я пришлю ответ в этот чат🕒
	`
)

func SupportHandler(ctx context.Context, event *common.Event, _ *gorm.DB) error {
	keyboard := common.NewInlineKeyboard()

	homeButton := common.NewButton("🏠 в начало", common.CallSetup().String())
	writeButton := common.NewButton("написать", common.CallWriteToSupport(event.GetMessage().GetChatIdStr()).String())
	keyboard.AppendAsStack(homeButton, writeButton)

	if _, err := event.EditCalbackMessage(
		ctx,
		fmt.Sprintf(SUPPORT_WELCOME_MESSAGE, "\n"),
		*keyboard.Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func WriteToSupportHandler(ctx context.Context, event *common.Event, _ *gorm.DB) error {
	event.SetNextHandler("send_to_support")

	if _, err := event.ReplyCallbackQuery(
		ctx,
		fmt.Sprintf(SUPPORT_INFO_MESSAGE, "\n\n"),
	); err != nil {
		return err
	}

	return nil
}

func SendToSupportHandler(ctx context.Context, event *common.Event, _ *gorm.DB) error {
	user := event.GetMessage().From

	userSupportRequest := event.GetMessage().Text

	if _, err := event.ReplyToSupport(
		ctx,
		fmt.Sprintf(
			"Новый запрос от пользователя\n"+
				"Username: `%s`\n"+
				"Name: `%s`\n"+
				"ChatId: `%s`\n"+
				"Request:\n\n`%s`",
			user.Username,
			user.FirstName,
			event.GetMessage().GetChatIdStr(),
			userSupportRequest,
		),
		telegram.WithMarkDown(),
	); err != nil {
		return err
	}

	if _, err := event.Reply(
		ctx,
		SUPPORT_THANKS_MESSAGE,
	); err != nil {
		return err
	}

	event.SetNextHandler("")

	return nil
}

func SendSupportResponseToUserHandler(ctx context.Context, event *common.Event, _ *gorm.DB) error {
	userRequest := event.GetMessage().ReplyToMessage.Text

	if userRequest == "" {
		event.Logger.Error("SendSupportResponseToUserHandler", "User request is empty", "skipping")
		return nil
	}

	userChatId, err := extractChatId(userRequest)
	if err != nil {
		event.Logger.Error("SendSupportResponseToUserHandler", "Error extracting chat id", err.Error())
		return err
	}

	messageFromSupport := event.GetMessage().Text
	message := fmt.Sprintf("Сообщение от команды поддержки🙌\n\n%s", messageFromSupport)

	if _, err := event.ReplyToUser(
		ctx,
		userChatId,
		message,
	); err != nil {
		return err
	}

	return nil
}

func extractChatId(text string) (string, error) {
	re := regexp.MustCompile(`ChatId:\s*(\d+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("ChatId not found")
}
