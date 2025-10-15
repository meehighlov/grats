package condition

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type Condition interface {
	Check(ctx context.Context, update *telegram.Update) (bool, error)
}
