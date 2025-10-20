package action

import (
	"context"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type Action func(ctx context.Context, update *telegram.Update) error
