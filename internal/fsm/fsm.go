package fsm

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/state"
	"github.com/meehighlov/grats/internal/fsm/store"
)

type FSM struct {
	states      []*state.State
	logger      *slog.Logger
	stateStore  store.StateStore
	middlewares []action.Action
}

func New(logger *slog.Logger, stateStore store.StateStore) *FSM {
	return &FSM{
		states:      []*state.State{},
		logger:      logger,
		stateStore:  stateStore,
		middlewares: []action.Action{},
	}
}
