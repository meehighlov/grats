package app

import (
	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/fsm"
	"github.com/meehighlov/grats/internal/infra/postgres"
	"github.com/meehighlov/grats/internal/infra/redis"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/server"
	"github.com/meehighlov/grats/internal/services"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	db := postgres.New(cfg, logger)
	tx := postgres.TransactionWrapper(cfg, db)

	redis := redis.New(cfg, logger)

	repositories := repositories.New(cfg, logger, db, redis)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	services := services.New(cfg, logger, repositories, clients, builders)

	fsm := fsm.New(logger, repositories.State)
	RegisterStates(fsm, services, cfg, clients, repositories, tx)

	server := server.New(cfg, logger, builders, clients.Telegram, fsm)
	server.Serve()
}
