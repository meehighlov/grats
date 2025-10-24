package clients

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
)

type Clients struct {
	Telegram *telegram.Client
}

func New(cfg *config.Config, logger *slog.Logger) *Clients {
	return &Clients{
		Telegram: telegram.New(cfg, logger),
	}
}

func (c *Clients) Close() error {
	return nil
}
