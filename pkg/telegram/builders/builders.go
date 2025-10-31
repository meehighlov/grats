package builders

import (
	callbackdata "github.com/meehighlov/grats/pkg/telegram/builders/callback_data"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	"github.com/meehighlov/grats/pkg/telegram/builders/pagination"
	"github.com/meehighlov/grats/pkg/telegram/config"
)

type Builders struct {
	CallbackData   *callbackdata.Builder
	InlineKeyboard *inlinekeyboard.Builder
	Pagination     *pagination.Builder
}

func New(cfg *config.Config) *Builders {
	return &Builders{
		CallbackData:   callbackdata.New(),
		InlineKeyboard: inlinekeyboard.New(),
		Pagination:     pagination.New(cfg, callbackdata.New(), inlinekeyboard.New()),
	}
}
