package pagination

import (
	callbackdata "github.com/meehighlov/grats/internal/builders/callback_data"
	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/config"
)

type Builder struct {
	callbackDataBuilder *callbackdata.Builder
	keyboardBuilder     *inlinekeyboard.Builder
	BaseOffset          int
	Limit               int
}

func New(cfg *config.Config, callbackDataBuilder *callbackdata.Builder, keyboardBuilder *inlinekeyboard.Builder) *Builder {
	return &Builder{
		callbackDataBuilder: callbackDataBuilder,
		keyboardBuilder:     keyboardBuilder,
		BaseOffset:          5,
		Limit:               cfg.ListLimitLen,
	}
}
