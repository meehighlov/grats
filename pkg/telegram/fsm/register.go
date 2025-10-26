package fsm

import (
	"log"
	"strconv"

	"github.com/meehighlov/grats/pkg/telegram/fsm/action"
	"github.com/meehighlov/grats/pkg/telegram/fsm/condition"
	"github.com/meehighlov/grats/pkg/telegram/fsm/state"
)

func (f *FSM) Activate(
	action action.Action,
	condition condition.Condition,
	opts ...state.StateOption,
) *state.State {
	stateId := strconv.Itoa(len(f.states) + 1)

	s := state.New(stateId, action, condition)

	for _, opt := range opts {
		opt(s)
	}

	if _, exists := f.states[s.GetID()]; exists {
		log.Fatalf("state with id %s already exists", s.GetID())
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
	ready := state.New("0", nil, nil)

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
