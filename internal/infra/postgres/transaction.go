package postgres

import (
	"context"
	"errors"

	"github.com/meehighlov/grats/internal/config"
	"gorm.io/gorm"
)

type DB struct {
	cfg *config.Config
	db  *gorm.DB
}

func TransactionWrapper(cfg *config.Config, db *gorm.DB) *DB {
	return &DB{
		cfg: cfg,
		db:  db,
	}
}

func (d *DB) Tx(ctx context.Context, fn func(context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, d.cfg.TxKey, tx)
		return fn(ctx)
	})
}

func (d *DB) GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(d.cfg.TxKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("not found transaction in context")
	}
	return tx, nil
}
