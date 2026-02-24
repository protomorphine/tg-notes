// Package log provides use full utiities for slog
package log

import (
	"log/slog"

	"github.com/google/uuid"
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

// ReqID creates a slog attribute for a request ID.
func ReqID(id uuid.UUID) slog.Attr {
	return slog.Attr{
		Key:   "reqID",
		Value: slog.StringValue(id.String()),
	}
}
