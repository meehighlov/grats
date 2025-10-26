package server

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
)

type Server struct {
	logger        *slog.Logger
	handleTimeout time.Duration
	telegram      *telegram.Client
	builders      *builders.Builders
	allowedUsers  []string
	webServer     *http.Server
	wgWebServer   sync.WaitGroup
	shutdownChan  chan struct{}
	cfg           *config.Config
	updateHandler UpdateHandler

	wgWorkerPool sync.WaitGroup
	workerCount  int
	workerCtx    context.Context
	workerCancel context.CancelFunc
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	builders *builders.Builders,
	telegram *telegram.Client,
	updateHandler UpdateHandler,
) *Server {
	return &Server{
		logger:        logger,
		builders:      builders,
		telegram:      telegram,
		handleTimeout: time.Duration(cfg.TelegramHandlerTimeoutSec) * time.Second,
		allowedUsers:  cfg.AdminList(),
		cfg:           cfg,
		shutdownChan:  make(chan struct{}),
		wgWebServer:   sync.WaitGroup{},
		wgWorkerPool:  sync.WaitGroup{},
		workerCount:   cfg.TelegramPollingWorkers,
		workerCtx:     context.Background(),
		workerCancel:  func() {},
		updateHandler: updateHandler,
	}
}
