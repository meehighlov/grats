package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/updates", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if !strings.EqualFold(token, s.cfg.TelegramWebhookToken) {
			s.logger.Warn("Invalid webhook token", "received", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var update models.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			s.logger.Error("Failed to decode update", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		s.logger.Debug("Received update", "update_id", update.UpdateId)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.TelegramHandlerTimeoutSec)*time.Second)
		defer cancel()

		if err := s.HandleUpdate(ctx, &update); err != nil {
			s.logger.Error("Failed to handle update", "error", err)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Health check requested")

		response := map[string]string{
			"status": "OK",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.logger.Error("Failed to encode health check response", "error", err)
		}
	})

	addr := s.cfg.TelegramWebhookAddress
	if s.cfg.TelegramUseTLS {
		addr = s.cfg.TelegramWebhookTLSAddress
	}
	s.webServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s.wgWebServer.Add(1)

	go func() {
		defer s.wgWebServer.Done()

		if s.cfg.TelegramUseTLS {
			s.logger.Info("Starting webhook server with TLS", "addr", s.webServer.Addr)
			if err := s.webServer.ListenAndServeTLS(s.cfg.TelegramWebhookTLSCertFile, s.cfg.TelegramWebhookTLSKeyFile); err != nil && err != http.ErrServerClosed {
				s.logger.Error("HTTP server error", "error", err)
			}
		} else {
			s.logger.Info("Starting webhook server without TLS", "addr", s.webServer.Addr)

			if err := s.webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Error("HTTP server error", "error", err)
			}
		}
	}()

	<-s.shutdownChan

	s.logger.Info("Webhook server stopping")
	return nil
}

func (s *Server) Stop() error {
	if s.webServer != nil {
		s.logger.Info("Stopping webhook server")

		close(s.shutdownChan)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.webServer.Shutdown(ctx); err != nil {
			s.logger.Error("Server shutdown error", "error", err)
		}

		s.wgWebServer.Wait()
		s.logger.Info("Webhook server stopped")
	}

	return nil
}

func (s *Server) Serve() error {
	if !s.cfg.TelegramUseWebook {
		return s.Polling()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go s.Start()

	<-signalChan
	return s.Stop()
}

func (s *Server) Polling() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates := s.telegram.GetUpdatesChannel(ctx)

	s.logger.Info("Starting polling mode with worker pool", "workers", s.workerCount)

	s.workerCtx, s.workerCancel = context.WithCancel(context.Background())

	s.logger.Info("Starting worker pool", "workers", s.workerCount)

	for i := range s.workerCount {
		s.wgWorkerPool.Add(1)
		go s.worker(i, &updates)
	}

	s.wgWorkerPool.Wait()

	return nil
}
