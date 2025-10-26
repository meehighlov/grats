package builders

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders/short_id"
	"github.com/meehighlov/grats/internal/config"
	callbackdata "github.com/meehighlov/grats/pkg/telegram/builders/callback_data"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	pagination "github.com/meehighlov/grats/pkg/telegram/builders/pagination"
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
		PaginationBuilder:   pagination.New(&cfg.Telegram, callbackDataBuilder, keyboardBuilder),
	}
}
