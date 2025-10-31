package clients

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
)

type Clients struct{}

func New(cfg *config.Config, logger *slog.Logger) *Clients {
	return &Clients{}
}
