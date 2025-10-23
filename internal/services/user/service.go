package user

import (
	"context"
	"log/slog"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/repositories"
)

type Common interface {
	RegisterOrUpdateUser(ctx context.Context, update *telegram.Update) error
}

type Service struct {
	common       Common
	logger       *slog.Logger
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
	common Common,
) *Service {
	return &Service{
		common:       common,
		logger:       logger,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		cfg:          cfg,
	}
}
