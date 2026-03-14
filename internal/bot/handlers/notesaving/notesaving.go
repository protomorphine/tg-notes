// Package notesaving provides default handler for bot
package notesaving

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"text/template"

	"protomorphine/tg-notes/internal/app/usecases/notesaving"
	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	successTemplate  = "resources/save_success.tmpl"
	errorTemplate    = "resources/save_err.tmpl"
	emptyMsgTemplate = "resources/empty_message.tmpl"
)

var (
	//go:embed resources
	templatesFS embed.FS

	templates map[string]*template.Template
)

func init() {
	templates = make(map[string]*template.Template)

	if tmpl, err := template.ParseFS(templatesFS, successTemplate); err != nil {
		templates[successTemplate] = tmpl
	}

	if tmpl, err := template.ParseFS(templatesFS, errorTemplate); err != nil {
		templates[errorTemplate] = tmpl
	}

	if tmpl, err := template.ParseFS(templatesFS, emptyMsgTemplate); err != nil {
		templates[emptyMsgTemplate] = tmpl
	}
}

// MessageSender is an interface for sending messages.
//
//mockery:generate: true
type MessageSender interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

// Handler represents the notesaving handler for the bot.
type Handler func(ctx context.Context, sender MessageSender, update *models.Update)

// New creates a new notesaving Handler.
func New(logger *slog.Logger, saver notesaving.NoteSaver) Handler {
	return func(ctx context.Context, sender MessageSender, update *models.Update) {
		const op = "bot.handlers.add"
		logger := logger.With(log.Op(op), log.ReqID(middleware.GetReqID(ctx)))

		if update.Message == nil {
			logger.Warn("nil message received")
			return
		}

		messageID := update.Message.ID
		chatID := update.Message.Chat.ID

		text := extractNoteText(update.Message)

		if text == "" {
			logger.Warn("received message with empty text and caption")

			if message, err := render(emptyMsgTemplate, struct{}{}); err == nil {
				sendMessage(ctx, logger, sender, chatID, messageID, message)
			} else {
				logger.Error("error while rendering template", log.Err(err))
			}

			return
		}

		res, err := saver.Save(ctx, text)
		if err != nil {
			logger.Error("error occured while saving new note", log.Err(err))

			if message, err := render(errorTemplate, struct{}{}); err == nil {
				sendMessage(ctx, logger, sender, chatID, messageID, message)
			} else {
				logger.Error("error while rendering template", log.Err(err))
			}

			return
		}

		logger.Info("new note saved")

		if message, err := render(successTemplate, res); err == nil {
			sendMessage(ctx, logger, sender, chatID, messageID, message)
		} else {
			logger.Error("error while rendering template", log.Err(err))
		}
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

func extractNoteText(message *models.Message) string {
	text := message.Text

	if text == "" {
		text = message.Caption
	}

	return text
}

func render(templatePath string, args any) (string, error) {
	tmpl, ok := templates[templatePath]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", templatePath)
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}
