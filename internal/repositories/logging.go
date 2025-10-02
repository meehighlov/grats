package repositories

import (
	"log"
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type appLoggerWrapper struct {
	logger *slog.Logger
}

func WrapAppLogger(logger *slog.Logger) gormlogger.Interface {
	appLogger := gormlogger.New(
		log.New(&appLoggerWrapper{logger: logger}, "", 0),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	return appLogger
}

func (l *appLoggerWrapper) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p))
	return len(p), nil
}
