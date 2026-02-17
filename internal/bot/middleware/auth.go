package middleware

import (
	"context"
	_ "embed"
	"log/slog"

	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

//go:embed resources/auth_err.tmpl
var authErrMsg string

// NewAuth function creates middleware to authorize user requests by telegram ID.
func NewAuth(logger *slog.Logger, cfg *config.BotConfig) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		logger := logger.With(slog.String("component", "middleware/auth"))

		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger := logger.With(slog.String("reqID", GetReqID(ctx).String()))

			if update.Message == nil {
				logger.Warn("got nil message in update")
				return
			}

			if update.Message.From.ID == cfg.AllowedUserID {
				logger.Info("successfully authorized new request")

				next(ctx, b, update)
				return
			}

			logger.Warn("sender ID missmatch allowed user ID", slog.Int64("fromID", update.Message.From.ID))

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      authErrMsg,
				ParseMode: models.ParseModeMarkdown,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			if err != nil {
				logger.Error("error while sending message", log.Err(err))
			}
		}
	}
}
