package app

import (
	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/orchestrators"
	"github.com/meehighlov/grats/internal/pagination"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/server"
	"github.com/meehighlov/grats/internal/services"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	db := db.New(cfg, logger)
	repositories := repositories.New(cfg, logger)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	constants := constants.New(cfg)
	pagination := pagination.New(cfg, builders)
	services := services.New(cfg, logger, repositories, clients, builders, constants, pagination)
	orchestrators := orchestrators.New(cfg, logger, db, services)

	server := server.New(cfg, logger, orchestrators, clients, constants, builders)
	server.Serve()
}
