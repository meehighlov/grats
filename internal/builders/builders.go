package builders

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders/short_id"
	"github.com/meehighlov/grats/internal/config"
)

type Builders struct {
	ShortIdBuilder *short_id.Builder
}

func New(cfg *config.Config, logger *slog.Logger) *Builders {
	return &Builders{
		ShortIdBuilder: short_id.New(cfg),
	}
}
