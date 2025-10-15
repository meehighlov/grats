package handler

import (
	"context"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type HandlerType func(ctx context.Context, update *telegram.Update) error
