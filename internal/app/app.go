package app

import (
	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
	"github.com/meehighlov/grats/internal/infra/redis"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/services"
	"github.com/meehighlov/grats/pkg/telegram/fsm"
	"github.com/meehighlov/grats/pkg/telegram/server"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	db := postgres.New(cfg, logger)

	redis := redis.New(cfg, logger)

	repositories := repositories.New(cfg, logger, db, redis)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	services := services.New(cfg, logger, repositories, clients, builders, db)

	fsm := fsm.New(logger, repositories.State)
	RegisterStates(fsm, services, cfg, clients, repositories, db)

	server := server.New(&cfg.Telegram, logger, fsm)
	server.Serve()
}
