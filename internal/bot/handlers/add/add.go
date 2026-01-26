// Package add provides handler for /add tg bot command
package add

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"protomorphine/tg-notes/internal/bot/middleware"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	saveErrMsg      = "unable to save new note :\\("
	saveSuccessMsg  = "note saved successfully"
	emptyMessageMsg = "can't process empty message"

	Cmd = "add"
)

type NoteAdder interface {
	Add(title, text string) error
}

func New(logger *slog.Logger, adder NoteAdder) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "bot.handlers.add"
		log := logger.With(sl.Op(op), slog.String("reqID", middleware.GetReqID(ctx).String()))

		if update.Message == nil {
			log.Error("empty message received")
			return
		}

		chatID := update.Message.Chat.ID

		text := strings.TrimPrefix(update.Message.Text, "/"+Cmd)
		text = strings.TrimSpace(text)

		if text == "" {
			log.Warn("empty message received")
			sendMessage(ctx, log, b, chatID, emptyMessageMsg)

			return
		}

		title := fmt.Sprintf("tg-notes bot %v", time.Now().Format(time.DateTime))

		if err := adder.Add(title, text); err != nil {
			log.Error("error occured while saving new note", sl.Err(err))
			sendMessage(ctx, log, b, chatID, saveErrMsg)

			return
		}

		log.Info("new note saved")
		sendMessage(ctx, log, b, chatID, saveSuccessMsg)
	}
}

func sendMessage(
	ctx context.Context,
	log *slog.Logger,
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
		log.Error("error occured while sending message", sl.Err(err))
	}
}
