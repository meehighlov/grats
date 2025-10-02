package builders

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders/callback_data"
	"github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/builders/short_id"
	"github.com/meehighlov/grats/internal/config"
)

type Builders struct {
	ShortIdBuilder      *short_id.Builder
	CallbackDataBuilder *callbackdata.Builder
	KeyboardBuilder     *inlinekeyboard.Builder
}

func New(cfg *config.Config, logger *slog.Logger) *Builders {
	return &Builders{
		ShortIdBuilder:      short_id.New(cfg),
		CallbackDataBuilder: callbackdata.New(),
		KeyboardBuilder:     inlinekeyboard.New(),
	}
}
