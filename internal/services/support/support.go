package support

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/pkg/telegram"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) Support(ctx context.Context, scope *telegram.Scope) error {
	keyboard := scope.Keyboard()

	callbackDataBuilder := scope.CallbackData()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_WRITE, callbackDataBuilder.Build("", s.cfg.Constants.CMD_SUPPORT_WRITE, "").String()),
		keyboard.NewButton(s.cfg.Constants.BTN_CANCEL, callbackDataBuilder.Build("", s.cfg.Constants.CMD_SUPPORT_CANCEL, "").String()),
	)

	if _, err := scope.Reply(ctx, s.cfg.Constants.SUPPORT_REQUEST_MESSAGE, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) SupportWrite(ctx context.Context, scope *telegram.Scope) error {
	keyboard := scope.Keyboard()
	callbackData := scope.CallbackData()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_CANCEL, callbackData.Build("", s.cfg.Constants.CMD_SUPPORT_CANCEL, "").String()),
	)

	if _, err := scope.Edit(ctx, s.cfg.Constants.SUPPORT_SEND_MESSAGE, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelSupportCall(ctx context.Context, scope *telegram.Scope) error {
	s.repositories.Cache.Reset(ctx, scope.Update().GetChatIdStr())

	if err := scope.DeleteMessage(ctx, scope.Update().GetChatIdStr(), strconv.Itoa(scope.Update().CallbackQuery.Message.MessageId)); err != nil {
		s.logger.Error("Failed to delete message", "error", err)
		return err
	}

	return nil
}

func (s *Service) SendSupportMessage(ctx context.Context, scope *telegram.Scope) error {
	message := scope.Update().GetMessage()

	if len(message.Text) > 2000 {
		scope.Reply(ctx, s.cfg.Constants.SUPPORT_MESSAGE_TOO_LONG)
		return nil
	}

	username := message.From.Username
	if username == "" {
		username = message.From.FirstName
	}

	supportMessage := fmt.Sprintf(
		s.cfg.Constants.SUPPORT_MESSAGE_TEMPLATE,
		message.GetChatIdStr(),
		username,
		strconv.Itoa(message.From.Id),
		message.Text,
	)

	if _, err := scope.ReplyTo(ctx, supportMessage, s.cfg.SupportChatId); err != nil {
		s.logger.Error("Failed to send support message", "error", err)
		scope.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE)
		return err
	}

	keyboard := scope.Keyboard()
	if _, err := scope.Reply(ctx, s.cfg.Constants.SUPPORT_MESSAGE_SENT, tgc.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) ProcessSupportReply(ctx context.Context, scope *telegram.Scope) error {
	message := scope.Update().GetMessage()

	if message.GetChatIdStr() != s.cfg.SupportChatId {
		return nil
	}

	if !message.IsReply() {
		return nil
	}

	chatId := s.parseChatIdFromMessage(message.ReplyToMessage.Text)
	if chatId == "" {
		s.logger.Error("Failed to parse chatid from support reply")
		return nil
	}

	replyMessage := fmt.Sprintf(s.cfg.Constants.SUPPORT_REPLY_TEMPLATE, message.Text)

	if _, err := scope.ReplyTo(ctx, replyMessage, chatId); err != nil {
		s.logger.Error("Failed to send reply to user", "error", err, "chatId", chatId)
		return err
	}

	s.logger.Info("Support reply sent successfully", "chatId", chatId)

	return nil
}

func (s *Service) parseChatIdFromMessage(text string) string {
	chatIdPrefix := s.cfg.Constants.SUPPORT_CHAT_ID_PREFIX
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return ""
	}

	chatId := strings.TrimPrefix(lines[0], chatIdPrefix)

	return chatId
}
