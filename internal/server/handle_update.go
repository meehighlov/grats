package server

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type UpdateHandler interface {
	Handle(ctx context.Context, update *telegram.Update) error
}

func (s *Server) HandleUpdate(ctx context.Context, update *telegram.Update) error {
	// TODO recover from panic
	// TODO call answerCallbackQuery
	s.logger.Info("Start handling", "update", update)
	err := s.updateHandler.Handle(ctx, update)
	if err != nil {
		s.logger.Error("Error", "update", update, "error", err)
	} else {
		s.logger.Info("Success", "update", update)
	}
	return err
}
