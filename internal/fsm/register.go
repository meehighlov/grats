package fsm

import (
	"log"

	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/condition"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func (f *FSM) Activate(
	action action.Action,
	condition condition.Condition,
	opts ...state.StateOption,
) {
	s := state.New(action, condition)

	if _, exists := f.states[s.GetID()]; exists {
		log.Fatalf("state with id %s already registered", s.GetID())
	}

	for _, opt := range opts {
		opt(s)
	}

	f.states[s.GetID()] = s
}
