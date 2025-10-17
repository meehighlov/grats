package state

import (
	"github.com/google/uuid"
	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
)

type initialState string

const (
	READY initialState = "ready"
)

func (s initialState) String() string {
	return string(s)
}

type StateOption func(*State) error

type InputState struct {
	FromStateId string
}

type OutputState struct {
	ActionError error
	ToStateId string
}

type State struct {
	id           string
	beforeAction []action.Action
	action       action.Action
	condition    condition.Condition

	inputStates  map[string]*InputState
	outputStates map[string]*OutputState
}

func New(action action.Action, condition condition.Condition) *State {
	return &State{
		id:           uuid.NewString(),
		beforeAction: nil,
		action:       action,
		condition:    condition,
		inputStates:  make(map[string]*InputState),
		outputStates: make(map[string]*OutputState),
	}
}
