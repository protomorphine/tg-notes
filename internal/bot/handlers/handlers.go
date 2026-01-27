// Package defaultHandler provides default handler for bot
package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	saveErrMsg      = "unable to save new note :\\("
	saveSuccessMsg  = "note saved successfully"
	emptyMessageMsg = "can't process empty message"
)

type NoteAdder interface {
	Add(title, text string) error
}

func NewDefault(logger *slog.Logger, adder NoteAdder) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.add"
		logger := logger.With(log.Op(op), slog.String("reqID", middleware.GetReqID(ctx).String()))

		if update.Message == nil {
			logger.Error("empty message received")
			return
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		if text == "" {
			logger.Warn("empty message received")
			sendMessage(ctx, logger, b, chatID, emptyMessageMsg)

			return
		}

		title := fmt.Sprintf("tg-notes bot %v", time.Now().Format(time.DateTime))

		if err := adder.Add(title, text); err != nil {
			logger.Error("error occured while saving new note", log.Err(err))
			sendMessage(ctx, logger, b, chatID, saveErrMsg)

			return
		}

		logger.Info("new note saved")
		sendMessage(ctx, logger, b, chatID, saveSuccessMsg)
	}
}

func sendMessage(
	ctx context.Context,
	logger *slog.Logger,
	b *bot.Bot,
	chatID int64,
	text string,
) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		logger.Error("error occured while sending message", log.Err(err))
	}
}
