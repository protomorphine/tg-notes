package log

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(err),
	}
}

func Op(op string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(op),
	}
}
