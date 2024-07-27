package lib

import (
	"log"
	"log/slog"

	"github.com/natefinch/lumberjack"
)

func MustSetupLogging(fileName string, useAsDefault bool, env string) *slog.Logger {
	envToLogLevel := map[string]slog.Level{}
	envToLogLevel["local"] = slog.LevelDebug
	envToLogLevel["prod"] = slog.LevelInfo

	logLevel, found := envToLogLevel[env]

	if !found {
		log.Fatalf("Error setting up logging, unknown env: %s", env)
	}

	// writes to file with file retention policy
	fileLogger := lumberjack.Logger{
		Filename:   fileName,
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

	if useAsDefault {
		slog.SetDefault(logger)
	}

	return logger
}
