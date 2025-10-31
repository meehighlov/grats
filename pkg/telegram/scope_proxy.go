package telegram

import (
	"context"

	callbackdata "github.com/meehighlov/grats/pkg/telegram/builders/callback_data"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	"github.com/meehighlov/grats/pkg/telegram/builders/pagination"
	"github.com/meehighlov/grats/pkg/telegram/client"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

func (s *Scope) Update() *models.Update {
	return s.update
}

func (s *Scope) UserID() string {
	return s.update.GetChatIdStr()
}

func (s *Scope) Reply(ctx context.Context, message string, opts ...client.SendMessageOption) (*models.Message, error) {
	return s.client.Reply(ctx, message, s.update, opts...)
}

func (s *Scope) ReplyTo(ctx context.Context, message string, chatId string, opts ...client.SendMessageOption) (*models.Message, error) {
	return s.client.SendMessage(ctx, chatId, message, opts...)
}

func (s *Scope) Edit(ctx context.Context, message string, opts ...client.SendMessageOption) (*models.Message, error) {
	return s.client.Edit(ctx, message, s.update, opts...)
}

func (s *Scope) SendFile(ctx context.Context, file []byte, filename string, opts ...client.SendMessageOption) (*models.SendDocumentResponse, error) {
	return s.client.SendFile(ctx, s.update, file, filename, opts...)
}

func (s *Scope) DeleteMessage(ctx context.Context, chatId string, messageId string) error {
	return s.client.DeleteMessage(ctx, chatId, messageId)
}

func (s *Scope) AnswerCallbackQuery(ctx context.Context) error {
	return s.client.AnswerCallbackQuery(ctx, s.update)
}

func (s *Scope) GetChatMember(ctx context.Context, chatId string) (*models.SingleChatMemberResponse, error) {
	return s.client.GetChatMember(ctx, chatId)
}

func (s *Scope) Keyboard() *inlinekeyboard.Builder {
	return s.builders.InlineKeyboard.NewKeyboard()
}

func (s *Scope) CallbackData() *callbackdata.Builder {
	return s.builders.CallbackData
}

func (s *Scope) Pagination() *pagination.Builder {
	return s.builders.Pagination
}
