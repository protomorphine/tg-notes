package main

import (
	"context"
	"fmt"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/handlers"
	"protomorphine/tg-notes/internal/bot/handlers/fallback"
	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/config"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
)

type webhookRemoveFunc func()

func newBot(logger *slog.Logger, cfg *config.BotConfig) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithErrorsHandler(handlers.NewErrorHandler(logger)),
		bot.WithDefaultHandler(fallback.New(logger)),
		bot.WithCheckInitTimeout(cfg.InitTimeout),
		bot.WithMiddlewares(
			middleware.NewReqID(),
			middleware.NewLog(logger),
		),
	}

	if logger.Enabled(context.Background(), slog.LevelDebug) {
		opts = append(opts,
			bot.WithDebug(),
			bot.WithDebugHandler(handlers.NewDebugHandler(logger)),
		)
	}

	return bot.New(cfg.Key, opts...)
}

func mustSetWebhook(ctx context.Context, logger *slog.Logger, b *bot.Bot, webhookURL string) webhookRemoveFunc {
	_, err := b.SetWebhook(ctx, &bot.SetWebhookParams{URL: webhookURL})
	if err != nil {
		panic(fmt.Errorf("error while setting webhook: %v", err))
	}

	return func() {
		_, err := b.DeleteWebhook(context.Background(), &bot.DeleteWebhookParams{DropPendingUpdates: true})
		if err != nil {
			logger.Error("error while deleting webhook", sl.Err(err))
			return
		}

		logger.Info("webhook was deleted")
	}
}
