// Package sl provides helper functions for "log/slog" 
package sl

import (
	"context"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/middleware"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ReqID(ctx context.Context) slog.Attr {
	return slog.Attr{
		Key:   "reqID",
		Value: slog.StringValue(middleware.GetReqID(ctx).String()),
	}
}

func Op(op string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(op),
	}
}
