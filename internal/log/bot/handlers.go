// Package bot provides helpers and logging handlers
package bot

import (
	"fmt"
	"log/slog"

	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
)

const component string = "go-telegram/bot"

// NewDebugHandler returns new logging handler for Debug level
func NewDebugHandler(logger *slog.Logger) bot.DebugHandler {
	return func(format string, args ...any) {
		logger.
			With(slog.String("compotent", component)).
			Debug(fmt.Sprintf(format, args...))
	}
}

// NewErrorHandler returns new logging handler for Error level
func NewErrorHandler(logger *slog.Logger) bot.ErrorsHandler {
	return func(err error) {
		logger.
			With(slog.String("compotent", component)).
			Error("error occured", log.Err(err))
	}
}
