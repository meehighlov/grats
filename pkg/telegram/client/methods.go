package client

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

func (c *Client) Reply(
	ctx context.Context,
	text string,
	update *models.Update,
	opts ...SendMessageOption,
) (*models.Message, error) {
	msg, err := c.SendMessage(ctx, update.GetChatIdStr(), text, opts...)
	return msg, err
}

func (c *Client) Edit(
	ctx context.Context,
	text string,
	update *models.Update,
	opts ...SendMessageOption,
) (*models.Message, error) {
	msg, err := c.EditMessageText(
		ctx,
		update.GetChatIdStr(),
		update.GetMessageIdStr(),
		text,
		opts...,
	)
	return msg, err
}

func (c *Client) SendFile(ctx context.Context, update *models.Update, file []byte, filename string, opts ...SendMessageOption) (*models.SendDocumentResponse, error) {
	return c.SendDocument(ctx, update.GetChatIdStr(), file, filename, opts...)
}
