package state

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/fsm/action"
	"github.com/meehighlov/grats/pkg/telegram/fsm/condition"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

func (s *State) Activate(ctx context.Context, update *models.Update) error {
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
	allowActivationFromAnyState := len(s.transitions) == 0
	if allowActivationFromAnyState {
		return true
	}

	return s.GetID() == stateId
}

func (s *State) AddTransition(transition *State) {
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

func (s *State) SetID(stateId string) {
	s.id = stateId
}

func (s *State) GetID() string {
	return s.id
}

func (s *State) GetTransitions() []*State {
	return s.transitions
}

func (s *State) Done(err error, currentStateId string) string {
	if err != nil {
		return currentStateId
	}
	if len(s.transitions) == 0 {
		return READY.String()
	}
	return s.GetID()
}
