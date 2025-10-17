package fsm

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/fsm/action"
	"github.com/meehighlov/grats/internal/fsm/state"
	"github.com/meehighlov/grats/internal/fsm/store"
)

type switchStateMode string

const (
	WhenReady    switchStateMode = "ready"
	FromAnyState switchStateMode = "any"
)

type FSM struct {
	states      map[string]*state.State
	logger      *slog.Logger
	stateStore  store.StateStore
	middlewares []action.Action
	switchMode  switchStateMode
}

func New(logger *slog.Logger, stateStore store.StateStore, switchMode switchStateMode) *FSM {
	return &FSM{
		states:      make(map[string]*state.State),
		logger:      logger,
		stateStore:  stateStore,
		middlewares: []action.Action{},
		switchMode:  switchMode,
	}
}
