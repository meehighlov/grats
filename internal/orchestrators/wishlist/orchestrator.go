package wishlist

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/services"
)

type Orchestrator struct {
	logger   *slog.Logger
	cfg      *config.Config
	services *services.Services
}

func New(cfg *config.Config, logger *slog.Logger, services *services.Services) *Orchestrator {
	return &Orchestrator{
		logger:   logger,
		cfg:      cfg,
		services: services,
	}
}
