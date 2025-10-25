package user

import (
	"context"
	"log/slog"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
	"github.com/meehighlov/grats/internal/repositories"
)

type UserRegistration interface {
	RegisterOrUpdateUser(ctx context.Context, update *telegram.Update) error
}

type Service struct {
	userRegistration UserRegistration
	logger           *slog.Logger
	db               *postgres.DB
	repositories     *repositories.Repositories
	clients          *clients.Clients
	builders         *builders.Builders
	cfg              *config.Config
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	db *postgres.DB,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	userRegistration UserRegistration,
) *Service {
	return &Service{
		userRegistration: userRegistration,
		logger:           logger,
		db:               db,
		repositories:     repositories,
		clients:          clients,
		builders:         builders,
		cfg:              cfg,
	}
}
