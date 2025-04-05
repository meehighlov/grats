package db

import (
	"context"
	"log"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type slogWriter struct {
	logger *slog.Logger
}

func (sw *slogWriter) Write(p []byte) (n int, err error) {
	sw.logger.Info(string(p))
	return len(p), nil
}

func MustSetup(dsn string, lgr *slog.Logger, runMigrations bool) {
	var err error

	gormLogger := logger.New(
		log.New(&slogWriter{logger: lgr}, "", 0),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	if runMigrations {
		if err := RunMigrations(context.Background(), lgr); err != nil {
			log.Fatal("Migration error:", err)
		}
	}

	lgr.Info("Database is ready")
}

func GetDB() *gorm.DB {
	return db
}

func GetDBConnection() (*gorm.DB, error) {
	return db, nil
}
