package middleware

import (
	"context"
	"log/slog"

	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// NewRecover creates middleware to recover from panic's during request pipeline.
func NewRecover(logger *slog.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		logger := logger.With(slog.String("component", "middleware/recover"))

		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			requestLogger := logger.With(log.ReqID(GetReqID(ctx)))

			defer func() {
				if val := recover(); val != nil {
					requestLogger.Error("recovered", slog.Any("val", val))
				}
			}()

			next(ctx, b, update)
		}
	}
}
