package user

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
)

type Repository struct {
	logger *slog.Logger
	cfg    *config.Config
	tx     *postgres.Tx
}

func New(cfg *config.Config, logger *slog.Logger, tx *postgres.Tx) *Repository {
	return &Repository{
		logger: logger,
		cfg:    cfg,
		tx:     tx,
	}
}
