package state

import (
	"github.com/google/uuid"
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
)

const (
	READY = "ready"
	ANY   = "any"
)

type StateOption func(*State) error

type Transition struct {
	ActionError error
	StateId     string
}

type State struct {
	ID           string
	beforeAction []action.Action
	action       action.Action
	condition    condition.Condition

	activationOnlyAfterStates []string

	transitions []*Transition
}

func New(action action.Action, condition condition.Condition) *State {
	return &State{
		ID:           uuid.NewString(),
		beforeAction: nil,
		action:       action,
		condition:    condition,
		transitions:  []*Transition{},
	}
}
