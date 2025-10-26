package action

import (
	"context"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

type Action func(ctx context.Context, update *models.Update) error
