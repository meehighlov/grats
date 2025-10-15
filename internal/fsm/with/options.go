package with

import (
	"github.com/meehighlov/grats/internal/fsm/handler"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func Transition(err error, status string) state.StateOption {
	return func(s *state.State) error {
		s.AddTransition(&state.Transition{
			HandlerError: err,
			Status:       status,
		})
		return nil
	}
}

func BeforeHandler(beforeHandler handler.HandlerType) state.StateOption {
	return func(s *state.State) error {
		if beforeHandler != nil {
			s.SetBeforeHandler(beforeHandler)
		}
		return nil
	}
}

func AllowedActivationStatus(status string) state.StateOption {
	return func(s *state.State) error {
		s.SetAllowedActivationStatus(status)
		return nil
	}
}
