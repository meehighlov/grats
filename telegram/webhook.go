package telegram

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebhookServer struct {
	addr         string
	handler      UpdateHandler
	client       *Client
	logger       *slog.Logger
	shutdownChan chan struct{}
	secretToken  string
	server       *http.Server
	wg           sync.WaitGroup
}

func NewWebhookServer(addr string, token string, secretToken string, handler UpdateHandler, logger *slog.Logger) *WebhookServer {
	return &WebhookServer{
		addr:         addr,
		handler:      handler,
		client:       NewClient(token, logger),
		logger:       logger,
		shutdownChan: make(chan struct{}),
		secretToken:  secretToken,
	}
}

func (ws *WebhookServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/updates", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if !strings.EqualFold(token, ws.secretToken) {
			ws.logger.Warn("Invalid webhook token", "received", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var update Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			ws.logger.Error("Failed to decode update", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ws.logger.Debug("Received update", "update_id", update.UpdateId)

		if err := ws.handler(update, ws.client); err != nil {
			ws.logger.Error("Failed to handle update", "error", err)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ws.logger.Debug("Health check requested")

		response := map[string]string{
			"status": "OK",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			ws.logger.Error("Failed to encode health check response", "error", err)
		}
	})

	ws.server = &http.Server{
		Addr:    ws.addr,
		Handler: mux,
	}

	ws.wg.Add(1)

	go func() {
		defer ws.wg.Done()
		ws.logger.Info("Starting webhook server", "addr", ws.addr)
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ws.logger.Error("HTTP server error", "error", err)
		}
	}()

	<-ws.shutdownChan

	ws.logger.Info("Webhook server stopping")
	return nil
}

func (ws *WebhookServer) Stop() {
	if ws.server != nil {
		ws.logger.Info("Stopping webhook server")

		close(ws.shutdownChan)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := ws.server.Shutdown(ctx); err != nil {
			ws.logger.Error("Server shutdown error", "error", err)
		}

		ws.wg.Wait()
		ws.logger.Info("Webhook server stopped")
	}
}

func StartWebhook(addr string, token string, secretToken string, handler UpdateHandler, logger *slog.Logger) *WebhookServer {
	server := NewWebhookServer(addr, token, secretToken, handler, logger)
	go server.Start()
	return server
}
