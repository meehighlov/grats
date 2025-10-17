package state

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
)

func (s *State) DoTransition(err error) string {
	for _, outputState := range s.outputStates {
		if outputState.ActionError == err {
			return outputState.ToStateId
		}
	}

	return READY
}

func (s *State) Activate(ctx context.Context, update *telegram.Update) error {
	for _, beforeAction := range s.beforeAction {
		if err := beforeAction(ctx, update); err != nil {
			return err
		}
	}
	return s.action(ctx, update)
}

func (s *State) Condition() condition.Condition {
	return s.condition
}

func (s *State) IsActivationAllowed(stateId string) bool {
	allowActivationFromAnyState := len(s.inputStates) == 0
	if allowActivationFromAnyState {
		return true
	}

	for _, inputState := range s.inputStates {
		if inputState.FromStateId == stateId {
			return true
		}
	}

	return false
}

func (s *State) AddInputState(inputState *InputState) {
	s.inputStates = append(s.inputStates, inputState)
}

func (s *State) SetBeforeAction(beforeAction action.Action) {
	s.beforeAction = append(s.beforeAction, beforeAction)
}

func (s *State) SetAction(action action.Action) {
	s.action = action
}

func (s *State) SetCondition(condition condition.Condition) {
	s.condition = condition
}

func (s *State) AddOutputState(outputState *OutputState) {
	s.outputStates = append(s.outputStates, outputState)
}

func (s *State) SetID(stateId string) {
	s.id = stateId
}
