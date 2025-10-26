package server

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type UpdateHandler interface {
	Handle(ctx context.Context, update *models.Update) error
}

func (s *Server) HandleUpdate(ctx context.Context, update *models.Update) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			stack := debug.Stack()
			s.logger.Error(
				"Panic recovered in HandleUpdate",
				"panic", r,
				"update", update,
				"stack", string(stack),
			)
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	s.logger.Info("Start handling", "update", update)
	err = s.updateHandler.Handle(ctx, update)
	if err != nil {
		s.logger.Error("Error", "update", update, "error", err)
	} else {
		s.logger.Info("Success", "update", update)
	}
	return err
}
