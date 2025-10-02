package repositories

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/repositories/user"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/internal/repositories/wish_list"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repositories struct {
	User     *user.Repository
	Wish     *wish.Repository
	WishList *wish_list.Repository
}

func New(cfg *config.Config, logger *slog.Logger) *Repositories {
	db, err := gorm.Open(postgres.Open(cfg.PGDSN), &gorm.Config{
		Logger: WrapAppLogger(logger),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := RunMigrations(context.Background(), cfg, logger, db); err != nil {
		log.Fatal("Migration error:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	logger.Info("Database connection established")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := sqlDB.PingContext(ctx); err != nil {
				logger.Error("Database ping failed", "error", err)
			} else {
				logger.Debug("Database ping successful")
			}
			cancel()
		}
	}()

	return &Repositories{
		User:     user.New(cfg, logger, db),
		Wish:     wish.New(cfg, logger, db),
		WishList: wish_list.New(cfg, logger, db),
	}
}
