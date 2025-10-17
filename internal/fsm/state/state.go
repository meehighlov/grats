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

	inputStates  []*InputState
	outputStates []*OutputState
}

func New(action action.Action, condition condition.Condition) *State {
	return &State{
		id:           uuid.NewString(),
		beforeAction: nil,
		action:       action,
		condition:    condition,
		inputStates:  []*InputState{},
		outputStates: []*OutputState{},
	}
}
