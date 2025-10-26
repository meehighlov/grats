package fsm

import (
	"log/slog"

	"github.com/meehighlov/grats/pkg/telegram/fsm/action"
	"github.com/meehighlov/grats/pkg/telegram/fsm/state"
	"github.com/meehighlov/grats/pkg/telegram/fsm/store"
)

type FSM struct {
	states      map[string]*state.State
	logger      *slog.Logger
	stateStore  store.StateStore
	middlewares []action.Action
}

func New(logger *slog.Logger, stateStore store.StateStore) *FSM {
	return &FSM{
		states:      make(map[string]*state.State),
		logger:      logger,
		stateStore:  stateStore,
		middlewares: []action.Action{},
	}
}
