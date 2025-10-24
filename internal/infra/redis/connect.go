package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/meehighlov/grats/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis           *redis.Client
	ctx             context.Context
	logger          *slog.Logger
	CacheExpiration time.Duration
}

func New(cfg *config.Config, logger *slog.Logger) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	client := &Client{
		Redis:           rdb,
		ctx:             context.Background(),
		logger:          logger,
		CacheExpiration: time.Duration(cfg.ChatCacheExpirationMinutes) * time.Minute,
	}

	if err := client.Ping(); err != nil {
		logger.Error("Redis connection failed", "error", err)
	} else {
		logger.Info("Redis connection successful", "addr", cfg.RedisAddr)
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := client.Ping(); err != nil {
				logger.Error("Redis ping failed", "error", err)
			} else {
				logger.Debug("Redis ping successful")
			}
		}
	}()

	return client
}

func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return c.Redis.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.Redis.Close()
}
