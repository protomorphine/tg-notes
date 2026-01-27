// Package log provides helpers and logging handlers
package log

import (
	"fmt"
	"log/slog"

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

		logger.Error("error occured", Err(err))
	}
}
