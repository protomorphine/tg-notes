package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NewLog(logger *slog.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		logger := logger.With(slog.String("component", "middleware/log"))

		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			entry := logger.With(
				slog.String("username", update.Message.From.Username),
				slog.Int64("chatID", update.Message.Chat.ID),
				slog.String("reqID", GetReqID(ctx).String()),
			)

			entry.Info("request accepted")

			t1 := time.Now()
			next(ctx, b, update)

			entry.Info("request completed", slog.String("duration", time.Since(t1).String()))
		}
	}
}
