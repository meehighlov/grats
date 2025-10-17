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

func InputState(stateId string) state.StateOption {
	return func(s *state.State) error {
		s.AddInputState(&state.InputState{
			FromStateId: stateId,
		})
		return nil
	}
}

func OutputState(err error, stateId string) state.StateOption {
	return func(s *state.State) error {
		s.AddOutputState(&state.OutputState{
			ActionError: err,
			ToStateId:     stateId,
		})
		return nil
	}
}

func Retry(err error) state.StateOption {
	return func(s *state.State) error {
		s.AddInputState(&state.InputState{
			FromStateId: s.GetID(),
		})
		s.AddOutputState(&state.OutputState{
			ActionError: err,
			ToStateId:   s.GetID(),
		})
		return nil
	}
}

func SuccessOutput(stateId string) state.StateOption {
	return func(s *state.State) error {
		s.AddOutputState(&state.OutputState{
			ActionError: nil,
			ToStateId:   stateId,
		})
		return nil
	}
}

func ID(stateId string) state.StateOption {
	return func(s *state.State) error {
		s.SetID(stateId)
		return nil
	}
}
