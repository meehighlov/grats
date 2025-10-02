package server

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/orchestrators"
)

type Server struct {
	logger        *slog.Logger
	handleTimeout time.Duration
	constants     *constants.Constants
	orchestrators *orchestrators.Orchestrators
	clients       *clients.Clients
	builders      *builders.Builders
	allowedUsers  []string
	webServer     *http.Server
	wgWebServer   sync.WaitGroup
	shutdownChan  chan struct{}
	cfg           *config.Config

	wgWorkerPool  sync.WaitGroup
	workerCount  int
	workerCtx    context.Context
	workerCancel context.CancelFunc
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	orchestrators *orchestrators.Orchestrators,
	clients *clients.Clients,
	constants *constants.Constants,
	builders *builders.Builders,
) *Server {
	return &Server{
		logger:        logger,
		orchestrators: orchestrators,
		clients:       clients,
		constants:     constants,
		builders:      builders,
		handleTimeout: time.Duration(cfg.TelegramHandlerTimeoutSec) * time.Second,
		allowedUsers:  cfg.AdminList(),
		cfg:           cfg,
		shutdownChan:  make(chan struct{}),
		wgWebServer:   sync.WaitGroup{},
		wgWorkerPool:  sync.WaitGroup{},
		workerCount:   cfg.TelegramPollingWorkers,
		workerCtx:     context.Background(),
		workerCancel:  func() {},
	}
}
