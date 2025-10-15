package server

import (
	"context"
	"runtime/debug"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type UpdateHandler interface {
	Handle(ctx context.Context, update *telegram.Update) error
}

func (s *Server) HandleUpdate(ctx context.Context, update *telegram.Update) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error(
				"Root handler",
				"recovered from panic, error", r,
				"stack", string(debug.Stack()),
				"update", update,
			)
			s.clients.Cache.Reset(ctx, update.GetChatIdStr())

			chatId := update.GetChatIdStr()
			if chatId != "" {
				s.clients.Telegram.Reply(ctx, s.constants.ERROR_MESSAGE, update)
				return
			}

			s.logger.Error(
				"Root handler",
				"recover from panic", "chatId was empty",
				"update", update,
			)
		}
	}()

	s.logger.Info("Start handling", "update", update)
	err := s.updateHandler.Handle(ctx, update)
	if err != nil {
		s.logger.Error("Error", "update", update, "error", err)
	} else {
		s.logger.Info("Success", "update", update)
	}
	return err
}
