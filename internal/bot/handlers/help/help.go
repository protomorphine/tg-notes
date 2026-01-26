// Package help contains handler for /help command
package help

import (
	"context"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/middleware"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	Cmd      = "help"
	helpText = `available commands:
\- */help*: provides some help information
\- */add*: add a new note
`
)

func New(logger *slog.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.help"
		logger := logger.With(
			sl.Op(op),
			slog.String("reqID", middleware.GetReqID(ctx).String()),
		)

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
