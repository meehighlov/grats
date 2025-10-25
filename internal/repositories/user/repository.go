package user

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
)

type Repository struct {
	logger *slog.Logger
	cfg    *config.Config
	db     *postgres.DB
}

func New(cfg *config.Config, logger *slog.Logger, db *postgres.DB) *Repository {
	return &Repository{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}
