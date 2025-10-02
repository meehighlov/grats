package telegram

import (
	"context"
)

func (c *Client) Reply(
	ctx context.Context,
	text string,
	update *Update,
	opts ...SendMessageOption,
) (*Message, error) {
	msg, err := c.SendMessage(ctx, update.GetChatIdStr(), text, opts...)
	return msg, err
}

func (c *Client) Edit(
	ctx context.Context,
	text string,
	update *Update,
	opts ...SendMessageOption,
) (*Message, error) {
	msg, err := c.EditMessageText(
		ctx,
		update.GetChatIdStr(),
		update.GetMessageIdStr(),
		text,
		opts...,
	)
	return msg, err
}

func (c *Client) SendFile(ctx context.Context, update *Update, file []byte, filename string, opts ...SendMessageOption) (*SendDocumentResponse, error) {
	return c.SendDocument(ctx, update.GetChatIdStr(), file, filename, opts...)
}
