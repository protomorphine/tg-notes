package main

import (
	"context"
	"fmt"
	"log/slog"

	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot"
)

type webhookRemoveFunc func()

func newBot(logger *slog.Logger, cfg *config.BotConfig, defaultHandler bot.HandlerFunc) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithErrorsHandler(log.NewErrorHandler(logger)),
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCheckInitTimeout(cfg.InitTimeout),
		bot.WithMiddlewares(
			middleware.NewRecover(logger),
			middleware.NewReqID(),
			middleware.NewAuth(logger, cfg),
			middleware.NewLog(logger),
		),
	}

	if logger.Enabled(context.Background(), slog.LevelDebug) {
		opts = append(opts,
			bot.WithDebug(),
			bot.WithDebugHandler(log.NewDebugHandler(logger)),
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
			logger.Error("error while deleting webhook", log.Err(err))
			return
		}

		logger.Info("webhook was deleted")
	}
}
