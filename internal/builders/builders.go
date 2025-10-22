package builders

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders/callback_data"
	"github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/builders/pagination"
	"github.com/meehighlov/grats/internal/builders/short_id"
	"github.com/meehighlov/grats/internal/config"
)

type Builders struct {
	ShortIdBuilder      *short_id.Builder
	CallbackDataBuilder *callbackdata.Builder
	KeyboardBuilder     *inlinekeyboard.Builder
	PaginationBuilder   *pagination.Builder
}

func New(cfg *config.Config, logger *slog.Logger) *Builders {
	callbackDataBuilder := callbackdata.New()
	keyboardBuilder := inlinekeyboard.New()
	return &Builders{
		ShortIdBuilder:      short_id.New(cfg),
		CallbackDataBuilder: callbackDataBuilder,
		KeyboardBuilder:     keyboardBuilder,
		PaginationBuilder:   pagination.New(cfg, callbackDataBuilder, keyboardBuilder),
	}
}
