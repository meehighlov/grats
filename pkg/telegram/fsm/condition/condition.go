package condition

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type Condition interface {
	Check(ctx context.Context, update *models.Update) (bool, error)
}
