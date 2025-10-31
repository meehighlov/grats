package fsm

import (
	"log"
	"strconv"
)

func (f *FSM) Activate(
	action Action,
	condition Condition,
	opts ...StateOption,
) *State {
	stateId := strconv.Itoa(len(f.states) + 1)

	s := NewState(stateId, action, condition)

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
	condition Condition,
	opts ...StateOption,
) *State {
	return f.Activate(f.reset, condition, opts...)
}

// only root states have to be reachable from READY state
func (f *FSM) setRootStates() {
	ready := NewState("0", nil, nil)

	delete(f.states, READY.String())

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

	f.states[READY.String()] = ready
}
