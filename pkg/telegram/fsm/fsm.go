package fsm

import (
	"log/slog"
)

type FSM struct {
	states      map[string]*State
	logger      *slog.Logger
	stateStore  StateStore
	middlewares []Action
}

func New(logger *slog.Logger, stateStore StateStore) *FSM {
	return &FSM{
		states:      make(map[string]*State),
		logger:      logger,
		stateStore:  stateStore,
		middlewares: []Action{},
	}
}
