package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"protomorphine/tg-notes/internal/bot/handlers/help"
	"protomorphine/tg-notes/internal/config"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
)

type CLIArgs struct {
	configPath string
}

func main() {
	args := mustParseAndValidateCLIArgs()

	cfg, err := config.Load(args.configPath)
	if err != nil {
		slog.Error("error occured while loading config", slog.Any("err", err))
		os.Exit(1)
	}

	logger := configureLogger(cfg.Environment, cfg.Logger)
	logger = logger.With(slog.String("env", cfg.Environment))

	logger.Info("starting tg-notes app")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	b, err := newBot(logger, &cfg.Bot)
	if err != nil {
		logger.Error("unable to initialize bot", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("successfully authorized in telegram api")

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, help.New(logger))

	removeWebhook := setWebhook(ctx, logger, b, cfg.Bot.WebHookURL)
	defer removeWebhook()

	server := &http.Server{
		Addr: cfg.HTTPServer.Addr,
		Handler: b.WebhookHandler(),
	}

	go func() {
		logger.Info("starting http server for receiving webhook's", slog.String("address", server.Addr))

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("error occured while running http server", sl.Err(err))
		}

		logger.Info("http server stopped")
	}()

	go func() {
		b.StartWebhook(ctx)
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("can't shutdown http server gracefully", sl.Err(err))
	}
}

func mustParseAndValidateCLIArgs() *CLIArgs {
	configPath := flag.String("config", "", "path to config file")

	flag.Parse()

	// check if config path is not empty and file exists
	if *configPath == "" {
		panic(errors.New("config path shouldn't be empty"))
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		panic(fmt.Errorf("file does not exists: %s", *configPath))
	}

	return &CLIArgs{configPath: *configPath}
}

