package state

import (
	"github.com/meehighlov/grats/internal/fsm/condition"
	"github.com/meehighlov/grats/internal/fsm/handler"
)

const (
	READY = "ready"
	ANY   = "any"
)

type StateOption func(*State) error

type Transition struct {
	HandlerError error
	Status       string
}

type State struct {
	beforeHandler           []handler.HandlerType
	handler                 handler.HandlerType
	condition               condition.Condition
	allowedActivationStatus string

	transitions []*Transition
}

func New(handler handler.HandlerType, condition condition.Condition) *State {
	return &State{
		handler:                 handler,
		condition:               condition,
		transitions:             []*Transition{},
		allowedActivationStatus: ANY,
	}
}
