package wish_list

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	logger *slog.Logger
	cfg *config.Config
}

func New(cfg *config.Config, logger *slog.Logger, db *gorm.DB) *Repository {
	return &Repository{
		db: db,
		logger: logger,
		cfg: cfg,
	}
}
