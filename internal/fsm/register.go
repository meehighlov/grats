package fsm

import (
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

	for _, opt := range opts {
		opt(s)
	}

	f.states = append(f.states, s)
}
