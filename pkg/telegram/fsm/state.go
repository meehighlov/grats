package fsm

import (
	"context"
)

type initialState string

const (
	READY initialState = ""
)

func (s initialState) String() string {
	return string(s)
}

type ActionData interface {
	UserID() string
}

type StateOption func(*State) error
type Action func(context.Context, ActionData) error

type Condition func(context.Context, ActionData) (bool, error)

type State struct {
	id           string
	beforeAction []Action
	action       Action
	condition    Condition

	transitions []*State
}

func NewState(id string, action Action, condition Condition) *State {
	return &State{
		id:           id,
		beforeAction: nil,
		action:       action,
		condition:    condition,
		transitions:  make([]*State, 0),
	}
}
