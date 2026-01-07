// Package fallback provides defautlt (fallback) tg bot handler
package fallback

import (
	"context"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/middleware"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const text = `didn't find any commands :\(
please use */help* to get information about available commands
`

func New(logger *slog.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.fallback"
		logger := logger.With(
			sl.Op(op),
			slog.String("reqId", middleware.GetReqID(ctx).String()),
		)

		logger.Info("got unknown command", slog.String("message", update.Message.Text))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      text,
			ParseMode: models.ParseModeMarkdown,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
		if err != nil {
			logger.Error("error while sending message", sl.Err(err))
		}
	}
}
