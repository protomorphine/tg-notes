package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// NewLog function    creates a middleware for log incoming requests.
func NewLog(logger *slog.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		logger := logger.With(slog.String("component", "middleware/log"))

		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger := logger.With(
				slog.String("username", update.Message.From.Username),
				slog.Int64("chatID", update.Message.Chat.ID),
				slog.String("reqID", GetReqID(ctx).String()),
			)

			t1 := time.Now()
			logger.Info("request accepted")

			next(ctx, b, update)

			logger.Info("request completed", slog.String("duration", time.Since(t1).String()))
		}
	}
}
