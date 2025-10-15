package support

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (s *Service) SupportHandler(ctx context.Context, update *telegram.Update) error {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.constants.BTN_WRITE, s.builders.CallbackDataBuilder.Build("", s.constants.CMD_SUPPORT_WRITE, "").String()),
		keyboard.NewButton(s.constants.BTN_CANCEL, s.builders.CallbackDataBuilder.Build("", s.constants.CMD_SUPPORT_CANCEL, "").String()),
	)

	if _, err := s.clients.Telegram.Reply(ctx, s.constants.SUPPORT_REQUEST_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) WriteHandler(ctx context.Context, update *telegram.Update) error {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.constants.BTN_CANCEL, s.builders.CallbackDataBuilder.Build("", s.constants.CMD_SUPPORT_CANCEL, "").String()),
	)

	if _, err := s.clients.Telegram.Edit(ctx, s.constants.SUPPORT_SEND_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelHandler(ctx context.Context, update *telegram.Update) error {
	s.clients.Cache.Reset(ctx, update.GetChatIdStr())

	if err := s.clients.Telegram.DeleteMessage(ctx, update.GetChatIdStr(), strconv.Itoa(update.CallbackQuery.Message.MessageId)); err != nil {
		s.logger.Error("Failed to delete message", "error", err)
		return err
	}

	return nil
}

func (s *Service) SendMessageHandler(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

	if len(message.Text) > 2000 {
		s.clients.Telegram.Reply(ctx, s.constants.SUPPORT_MESSAGE_TOO_LONG, update)
		return nil
	}

	username := message.From.Username
	if username == "" {
		username = message.From.FirstName
	}

	supportMessage := fmt.Sprintf(
		s.constants.SUPPORT_MESSAGE_TEMPLATE,
		message.GetChatIdStr(),
		username,
		strconv.Itoa(message.From.Id),
		message.Text,
	)

	if _, err := s.clients.Telegram.SendMessage(ctx, s.cfg.SupportChatId, supportMessage); err != nil {
		s.logger.Error("Failed to send support message", "error", err)
		s.clients.Telegram.Reply(ctx, s.constants.ERROR_MESSAGE, update)
		return err
	}

	keyboard := s.buildBackToMenuKeyboard()
	if _, err := s.clients.Telegram.Reply(ctx, s.constants.SUPPORT_MESSAGE_SENT, update, telegram.WithReplyMurkup(keyboard.Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildBackToMenuKeyboard() *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()
	return keyboard
}

func (s *Service) HandleSupportReply(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

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

	replyMessage := fmt.Sprintf(s.constants.SUPPORT_REPLY_TEMPLATE, message.Text)

	if _, err := s.clients.Telegram.SendMessage(ctx, chatId, replyMessage); err != nil {
		s.logger.Error("Failed to send reply to user", "error", err, "chatId", chatId)
		return err
	}

	s.logger.Info("Support reply sent successfully", "chatId", chatId)

	return nil
}

func (s *Service) parseChatIdFromMessage(text string) string {
	chatIdPrefix := s.constants.SUPPORT_CHAT_ID_PREFIX
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return ""
	}

	chatId := strings.TrimPrefix(lines[0], chatIdPrefix)

	return chatId
}
