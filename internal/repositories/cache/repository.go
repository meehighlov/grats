package cache

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/redis"
)

type Repository struct {
	logger *slog.Logger
	cfg    *config.Config
	redis  *redis.Client
}

func New(cfg *config.Config, logger *slog.Logger, redis *redis.Client) *Repository {
	return &Repository{
		logger: logger,
		cfg:    cfg,
		redis:  redis,
	}
}
