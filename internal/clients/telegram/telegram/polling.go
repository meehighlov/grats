package telegram

import (
	"context"
)

type UpdateHandler func(Update, *Client) error

func StartPolling(token string, handler UpdateHandler) error {
	withCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(token, nil)

	updates := client.GetUpdatesChannel(withCancel)

	for update := range updates {
		go handler(update, client)
	}

	return nil
}
