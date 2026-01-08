// Package handlers provides handlers to tg bot
package handlers

import (
	"fmt"
	"log/slog"

	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
)

const component string = "go-telegram/bot"

// NewDebugHandler function    returns new logging handler for Debug level
func NewDebugHandler(logger *slog.Logger) bot.DebugHandler {
	return func(format string, args ...any) {
		logger := logger.With(slog.String("compotent", component))

		logger.Debug(fmt.Sprintf(format, args...))
	}
}

// NewErrorHandler function    returns new logging handler for Error level
func NewErrorHandler(logger *slog.Logger) bot.ErrorsHandler {
	return func(err error) {
		logger := logger.With(slog.String("compotent", component))

		logger.Error("error occured", sl.Err(err))
	}
}
