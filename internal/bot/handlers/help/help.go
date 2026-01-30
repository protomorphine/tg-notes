// Package help contains handler for /help command
package help

import (
	"context"
	_ "embed"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const Cmd = "help"

//go:embed resources/help.tmpl
var helpText string

func New(logger *slog.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.help"
		logger := logger.With(
			log.Op(op),
			slog.String("reqID", middleware.GetReqID(ctx).String()),
		)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      helpText,
			ParseMode: models.ParseModeMarkdownV1,
		})
		if err != nil {
			logger.Error("error while sending message", log.Err(err))
		}
	}
}
