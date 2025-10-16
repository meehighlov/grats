package with

import (
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func Transition(err error, stateId string) state.StateOption {
	return func(s *state.State) error {
		s.AddTransition(&state.Transition{
			ActionError: err,
			StateId:     stateId,
		})
		return nil
	}
}

func BeforeAction(beforeAction action.Action) state.StateOption {
	return func(s *state.State) error {
		if beforeAction != nil {
			s.SetBeforeAction(beforeAction)
		}
		return nil
	}
}

func ActivationOnlyAfter(stateId string) state.StateOption {
	return func(s *state.State) error {
		s.ActivationOnlyAfter(stateId)
		return nil
	}
}

func ID(stateId string) state.StateOption {
	return func(s *state.State) error {
		s.ID = stateId
		return nil
	}
}
