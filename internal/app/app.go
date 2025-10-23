package app

import (
	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/fsm"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/server"
	"github.com/meehighlov/grats/internal/services"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	dbConn := db.New(cfg, logger)
	tx := db.TransactionWrapper(cfg, dbConn)

	repositories := repositories.New(cfg, logger)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	services := services.New(cfg, logger, repositories, clients, builders)

	fsm := fsm.New(logger, clients.Cache)
	RegisterStates(fsm, services, cfg, clients, tx)

	server := server.New(cfg, logger, clients, builders, fsm)
	server.Serve()
}
