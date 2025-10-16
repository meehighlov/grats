package state

import (
	"context"
	"slices"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
)

func (s *State) Next(err error) string {
	for _, transition := range s.transitions {
		if transition.ActionError == err {
			return transition.StateId
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
	fromAnyState := len(s.activationOnlyAfterStates) == 0
	if fromAnyState {
		return true
	}

	return slices.Contains(s.activationOnlyAfterStates, stateId)
}

func (s *State) AddTransition(transition *Transition) {
	s.transitions = append(s.transitions, transition)
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

func (s *State) ActivationOnlyAfter(stateId string) {
	s.activationOnlyAfterStates = append(s.activationOnlyAfterStates, stateId)
}
