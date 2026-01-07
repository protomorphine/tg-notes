package main

import (
	"log/slog"
	"os"
	"time"

	"protomorphine/tg-notes/internal/config"

	"github.com/lmittmann/tint"
)

const (
	// log levels
	Debug string = "DEBUG"
	Info  string = "INFO"
	Warn  string = "WARN"
	Error string = "ERROR"

	// application enviroments
	Local      string = "local"
	Production string = "prod"
)

func configureLogger(env string, logCfg config.LoggerConfig) *slog.Logger {
	var level slog.Level

	switch logCfg.MinLevel {
	case Info:
		level = slog.LevelInfo
	case Debug:
		level = slog.LevelDebug
	case Warn:
		level = slog.LevelWarn
	case Error:
		level = slog.LevelError
	}

	var handler slog.Handler

	switch env {
	case Local:
		handler = tint.NewHandler(
			os.Stdout,
			&tint.Options{
				AddSource:  true,
				Level:      level,
				TimeFormat: time.Kitchen,
			})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}
