package app

import (
	"log"
	"log/slog"

	"github.com/meehighlov/grats/internal/config"

	"github.com/natefinch/lumberjack"
)

func MustSetupLogging(cfg *config.Config) *slog.Logger {
	envToLogLevel := map[string]slog.Level{}
	envToLogLevel[config.LOCAL] = slog.LevelDebug
	envToLogLevel[config.PROD] = slog.LevelInfo

	logLevel, found := envToLogLevel[cfg.ENV]

	if !found {
		log.Fatalf("Error setting up logging, unknown env: %s", cfg.ENV)
	}

	// writes to file with file retention policy
	fileLogger := lumberjack.Logger{
		Filename:   cfg.LoggingFileName,
		MaxSize:    512, // megabytes
		MaxBackups: 1,
		MaxAge:     7, // days
		Compress:   false,
	}

	jsonHandler := slog.NewJSONHandler(
		&fileLogger,
		&slog.HandlerOptions{Level: logLevel},
	)

	logger := slog.New(jsonHandler)

	return logger
}
