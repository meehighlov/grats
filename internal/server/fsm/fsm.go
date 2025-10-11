package fsm

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
)

type FSM struct {
	nodes map[string]*node
	clients *clients.Clients
	constants *constants.Constants
	logger *slog.Logger
	cfg *config.Config
}

func New(cfg *config.Config, logger *slog.Logger, clients *clients.Clients, constants *constants.Constants) *FSM {
	return &FSM{
		nodes: map[string]*node{},
		clients: clients,
		logger: logger,
		cfg: cfg,
		constants: constants,
	}
}
