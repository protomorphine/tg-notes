package log

import (
	"context"
	"log/slog"
)

type discardHandler struct{}

// NewDiscardHandler creates a new slog.Handler that discards all log messages.
// This handler is primarily used for testing purposes.
func NewDiscardHandler() slog.Handler {
	return &discardHandler{}
}

func (h *discardHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *discardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h *discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *discardHandler) WithGroup(name string) slog.Handler {
	return h
}
