package middleware

import (
	"context"
	"log/slog"

	"protomorphine/tg-notes/internal/config"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// NewAuth function    creates a middleware to authorize user requests by telegram ID.
func NewAuth(logger *slog.Logger, cfg *config.BotConfig) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		logger := logger.With(slog.String("component", "middleware/auth"))

		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			requestLogger := logger.With(slog.String("reqID", GetReqID(ctx).String()))

			if update.Message.From.ID == cfg.AllowedUserID {
				requestLogger.Info("successfully authorized new request")

				next(ctx, b, update)
				return
			}

			requestLogger.Error("sender ID missmatch allowed user ID", slog.Int64("fromID", update.Message.From.ID))

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "you are not allowed to do this",
				ParseMode: models.ParseModeMarkdown,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			if err != nil {
				requestLogger.Error("error while sending message", sl.Err(err))
			}
		}
	}
}
