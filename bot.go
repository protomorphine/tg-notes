package main

import (
	"context"
	"fmt"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/handlers/help"
	"protomorphine/tg-notes/internal/bot/handlers/notesaving"
	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/log"
	botlog "protomorphine/tg-notes/internal/log/bot"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type webhookRemoveFunc func()

func newBot(logger *slog.Logger, cfg *config.BotConfig, defaultHandler notesaving.Handler) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithErrorsHandler(botlog.NewErrorHandler(logger)),
		bot.WithDefaultHandler(wrapHandler(defaultHandler)),
		bot.WithCheckInitTimeout(cfg.InitTimeout),
		bot.WithMiddlewares(
			middleware.NewReqID(),
			middleware.NewRecover(logger),
			middleware.NewAuth(logger, cfg),
			middleware.NewLog(logger),
		),
	}

	if logger.Enabled(context.Background(), slog.LevelDebug) {
		opts = append(opts,
			bot.WithDebug(),
			bot.WithDebugHandler(botlog.NewDebugHandler(logger)),
		)
	}

	b, err := bot.New(cfg.Key, opts...)
	if err != nil {
		return nil, err
	}

	// register additional command handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, help.Cmd, bot.MatchTypeCommand, help.New(logger))

	return b, nil
}

func wrapHandler(handler notesaving.Handler) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		handler(ctx, bot, update)
	}
}

func setWebhook(ctx context.Context, logger *slog.Logger, b *bot.Bot, webhookURL string) (webhookRemoveFunc, error) {
	_, err := b.SetWebhook(ctx, &bot.SetWebhookParams{URL: webhookURL})
	if err != nil {
		return nil, fmt.Errorf("set webhook error: %w", err)
	}

	return func() {
		_, err := b.DeleteWebhook(context.Background(), &bot.DeleteWebhookParams{DropPendingUpdates: true})
		if err != nil {
			logger.Error("error while deleting webhook", log.Err(err))
			return
		}

		logger.Info("webhook was deleted")
	}, nil
}
