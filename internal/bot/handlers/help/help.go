// Package help contains handler for /help command
package help

import (
	"context"
	"log/slog"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const helpText string = `available commands:
\- */help*: provides some help information
`

func New(logger *slog.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.help"
		logger := logger.With(sl.Op(op), sl.ReqID(ctx),)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      helpText,
			ParseMode: models.ParseModeMarkdown,
		})

		if err != nil {
			logger.Error("error while sending message", sl.Err(err))
		}
	}
}
