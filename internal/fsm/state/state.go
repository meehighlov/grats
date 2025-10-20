package state

import (
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
)

type initialState string

const (
	READY initialState = ""
)

func (s initialState) String() string {
	return string(s)
}

type StateOption func(*State) error

type State struct {
	id           string
	beforeAction []action.Action
	action       action.Action
	condition    condition.Condition

	transitions []*State
}

func New(id string, action action.Action, condition condition.Condition) *State {
	return &State{
		id:           id,
		beforeAction: nil,
		action:       action,
		condition:    condition,
		transitions:  make([]*State, 0),
	}
}
