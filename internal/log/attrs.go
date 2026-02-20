// Package log provides usefull utiities for slog
package log

import (
	"log/slog"
)

// Err creates a slog attribute for an error.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(err),
	}
}

// Op creates a slog attribute for an operation.
func Op(op string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(op),
	}
}
