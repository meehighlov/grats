package fsm

import (
	"context"
)

func (s *State) Activate(ctx context.Context, data ActionData) (err error) {
	for _, beforeAction := range s.beforeAction {
		if err = beforeAction(ctx, data); err != nil {
			return err
		}
	}
	return s.action(ctx, data)
}

func (s *State) Condition(ctx context.Context, data ActionData) (bool, error) {
	return s.condition(ctx, data)
}

func (s *State) AddTransition(transition *State) {
	s.transitions = append(s.transitions, transition)
}

func (s *State) SetBeforeAction(beforeAction Action) {
	s.beforeAction = append(s.beforeAction, beforeAction)
}

func (s *State) SetAction(action Action) {
	s.action = action
}

func (s *State) SetCondition(condition Condition) {
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
