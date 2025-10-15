package fsm

import (
	"github.com/meehighlov/grats/internal/fsm/condition"
	"github.com/meehighlov/grats/internal/fsm/handler"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func (f *FSM) Activate(
	handler handler.HandlerType,
	condition condition.Condition,
	opts ...state.StateOption,
) {
	s := state.New(handler, condition)

	for _, opt := range opts {
		opt(s)
	}

	f.states = append(f.states, s)
}
