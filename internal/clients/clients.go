package clients

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

type Clients struct {
	Telegram *tgc.Client
}

func New(cfg *config.Config, logger *slog.Logger) *Clients {
	return &Clients{
		Telegram: tgc.New(&cfg.Telegram, logger),
	}
}

func (c *Clients) Close() error {
	return nil
}
