// Package notesaving provides default handler for bot
package notesaving

import (
	"context"
	_ "embed"
	"log/slog"

	"protomorphine/tg-notes/internal/app/usecase/notesaving"
	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	//go:embed resources/save_err.tmpl
	saveErrMsg string

	//go:embed resources/save_success.tmpl
	saveSuccessMsg string

	//go:embed resources/empty_message.tmpl
	emptyMessageMsg string
)

//mockery:generate: true
type MessageSender interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

type Handler func(ctx context.Context, sender MessageSender, update *models.Update)

func New(logger *slog.Logger, saver *notesaving.Usecase) Handler {
	return func(ctx context.Context, sender MessageSender, update *models.Update) {
		const op = "bot.handlers.add"
		logger := logger.With(log.Op(op), slog.String("reqID", middleware.GetReqID(ctx).String()))

		if update.Message == nil {
			logger.Error("empty message received")
			return
		}

		messageID := update.Message.ID

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		if text == "" {
			text = update.Message.Caption

			if text == "" {
				logger.Warn("empty message received")
				sendMessage(ctx, logger, sender, chatID, messageID, emptyMessageMsg)

				return
			}
		}

		if err := saver.Save(ctx, text); err != nil {
			logger.Error("error occured while saving new note", log.Err(err))
			sendMessage(ctx, logger, sender, chatID, messageID, saveErrMsg)

			return
		}

		logger.Info("new note saved")
		sendMessage(ctx, logger, sender, chatID, messageID, saveSuccessMsg)
	}
}

func sendMessage(
	ctx context.Context,
	logger *slog.Logger,
	sender MessageSender,
	chatID int64,
	replyID int,
	text string,
) {
	_, err := sender.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyID,
		},
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		logger.Error("error occured while sending message", log.Err(err))
	}
}
