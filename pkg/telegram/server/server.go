package server

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/meehighlov/grats/pkg/telegram/client"
	"github.com/meehighlov/grats/pkg/telegram/config"
)

type Server struct {
	logger        *slog.Logger
	handleTimeout time.Duration
	telegram      *client.Client
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
	updateHandler UpdateHandler,
) *Server {
	telegram := client.New(cfg, logger)
	return &Server{
		logger:        logger,
		telegram:      telegram,
		handleTimeout: time.Duration(cfg.TelegramHandlerTimeoutSec) * time.Second,
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
