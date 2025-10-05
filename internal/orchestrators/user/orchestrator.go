package user

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/services"
	"gorm.io/gorm"
)

type Orchestrator struct {
	logger   *slog.Logger
	cfg      *config.Config
	services *services.Services
	db       *gorm.DB
}

func New(cfg *config.Config, logger *slog.Logger, services *services.Services, db *gorm.DB) *Orchestrator {
	return &Orchestrator{
		logger:   logger,
		cfg:      cfg,
		services: services,
		db:       db,
	}
}
