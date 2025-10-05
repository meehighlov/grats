package wish_list

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
)

type Repository struct {
	logger *slog.Logger
	cfg    *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) *Repository {
	return &Repository{
		logger: logger,
		cfg:    cfg,
	}
}
