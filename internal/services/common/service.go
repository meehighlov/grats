package common

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
	"github.com/meehighlov/grats/internal/repositories"
)

type Service struct {
	logger       *slog.Logger
	tx           *postgres.Tx
	repositories *repositories.Repositories
	clients      *clients.Clients
	builders     *builders.Builders
	cfg          *config.Config
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	tx *postgres.Tx,
) *Service {
	return &Service{
		logger:       logger,
		tx:           tx,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		cfg:          cfg,
	}
}
