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
) *state.State {
	s := state.New(action, condition)

	for _, opt := range opts {
		opt(s)
	}

	f.states[s.GetID()] = s

	f.setRootStates()

	return s
}

func (f *FSM) Reset(
	condition condition.Condition,
	opts ...state.StateOption,
) *state.State {
	return f.Activate(f.reset, condition, opts...)
}

// only root states have to be reachable from READY state
func (f *FSM) setRootStates() {
	ready := state.New(nil, nil)

	delete(f.states, state.READY.String())

	statesWithIncomingEdges := make(map[string]bool)

	for _, st := range f.states {
		for _, transition := range st.GetTransitions() {
			statesWithIncomingEdges[transition.GetID()] = true
		}
	}

	for stateId, st := range f.states {
		if !statesWithIncomingEdges[stateId] {
			ready.AddTransition(st)
		}
	}

	f.states[state.READY.String()] = ready
}
