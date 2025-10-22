package user

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (s *Service) Start(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

	err := s.RegisterOrUpdateUser(ctx, update)
	if err != nil {
		return err
	}

	username := message.From.Username
	if username == "" {
		username = message.From.FirstName
		if username == "" {
			username = s.cfg.Constants.GREETING_FRIEND
		}
	}

	hello := fmt.Sprintf(
		s.cfg.Constants.GREETING_TEMPLATE,
		username,
	)

	if _, err := s.clients.Telegram.Reply(ctx, hello, update); err != nil {
		return err
	}

	return nil
}
