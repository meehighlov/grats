package wishlist

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
	db           *postgres.DB
	repositories *repositories.Repositories
	clients      *clients.Clients
	builders     *builders.Builders
	cfg          *config.Config
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	db *postgres.DB,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
) *Service {
	return &Service{
		logger:       logger,
		db:           db,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		cfg:          cfg,
	}
}
