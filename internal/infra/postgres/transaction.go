package postgres

import (
	"context"
	"errors"

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

func (t *Tx) Atomic(ctx context.Context, fn func(context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, t.cfg.TxKey, tx)
		return fn(ctx)
	})
}

func (t *Tx) GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(t.cfg.TxKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("not found transaction in context")
	}
	return tx, nil
}
