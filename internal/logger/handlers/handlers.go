package handlers

import (
	"fmt"
	"log/slog"
	sl "protomorphine/tg-notes/internal/logger"
)

func NewDebugHandler(logger *slog.Logger) func(format string, args ...any) {
	return func(format string, args ...any) {
		logger := logger.With(slog.String("compotent", "tg-bot"))

		logger.Debug(fmt.Sprintf(format, args...))
	}
}

func NewErrorHandler(logger *slog.Logger) func(err error) {
	return func(err error) {
		logger := logger.With(slog.String("compotent", "tg-bot"))

		logger.Error("error occured", sl.Err(err))
	}
}
