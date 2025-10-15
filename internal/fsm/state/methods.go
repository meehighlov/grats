package state

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/fsm/condition"
	"github.com/meehighlov/grats/internal/fsm/handler"
)

func (s *State) Next(err error) string {
	for _, transition := range s.transitions {
		if transition.HandlerError == err {
			return transition.Status
		}
	}

	// no transition found - flow is done,
	// machine is ready
	return READY
}

func (s *State) Activate(ctx context.Context, update *telegram.Update) error {
	for _, beforeHandler := range s.beforeHandler {
		if err := beforeHandler(ctx, update); err != nil {
			return err
		}
	}
	return s.handler(ctx, update)
}

func (s *State) Condition() condition.Condition {
	return s.condition
}

func (s *State) IsActivationAllowed(status string) bool {
	return s.allowedActivationStatus == status || s.allowedActivationStatus == ANY
}

func (s *State) AddTransition(transition *Transition) {
	s.transitions = append(s.transitions, transition)
}

func (s *State) SetBeforeHandler(beforeHandler handler.HandlerType) {
	s.beforeHandler = append(s.beforeHandler, beforeHandler)
}

func (s *State) SetHandler(handler handler.HandlerType) {
	s.handler = handler
}

func (s *State) SetCondition(condition condition.Condition) {
	s.condition = condition
}

func (s *State) SetAllowedActivationStatus(status string) {
	s.allowedActivationStatus = status
}
