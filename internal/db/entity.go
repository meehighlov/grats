package db

import (
	"context"

	"gorm.io/gorm"
)

type Entity interface {
	GetId() string
	GreaterThan(other Entity) bool
	ButtonText() string
	Search(ctx context.Context, tx *gorm.DB, chatId, userId string) ([]Entity, error)
}
