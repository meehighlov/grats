package db

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"gorm.io/gorm"
)

type Tx struct {
	cfg *config.Config
	db  *gorm.DB
}

func TransactionWrapper(cfg *config.Config, db *gorm.DB) *Tx {
	return &Tx{
		cfg: cfg,
		db:  db,
	}
}

func (t *Tx) Wrap(fn func(ctx context.Context, update *telegram.Update) error) func(ctx context.Context, update *telegram.Update) error {
	return func(ctx context.Context, update *telegram.Update) error {
		return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, t.cfg.TxKey, tx)
			return fn(ctx, update)
		})
	}
}
