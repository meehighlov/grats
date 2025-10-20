package with

import (
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func BeforeAction(beforeAction action.Action) state.StateOption {
	return func(s *state.State) error {
		if beforeAction != nil {
			s.SetBeforeAction(beforeAction)
		}
		return nil
	}
}

func AcceptFrom(fromState *state.State) state.StateOption {
	return func(s *state.State) error {
		fromState.AddTransition(s)
		return nil
	}
}
