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

	"protomorphine/tg-notes/internal/bot/handlers"
	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/log"
	"protomorphine/tg-notes/internal/storage/git"
)

type CLIArgs struct {
	configPath string
}

func main() {
	args, err := parseAndValidateCLIArgs()
	if err != nil {
		slog.Error("error while parsing CLI args")
		os.Exit(1)
	}

	cfg, err := config.Load(args.configPath)
	if err != nil {
		slog.Error("error while loading config", slog.Any("err", err))
		os.Exit(1)
	}

	logger := configureLogger(cfg.Environment, &cfg.Logger)
	logger = logger.With(slog.String("env", cfg.Environment))

	logger.Info("starting tg-notes app")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	storage, err := git.New(&cfg.GitRepository)
	if err != nil {
		logger.Error("error while setting up storage", log.Err(err))
		os.Exit(1)
	}

	logger.Info("successfully initialized git storage")

	go storage.Processor(ctx, logger)

	b, err := newBot(logger, &cfg.Bot, handlers.NewDefault(logger, storage))
	if err != nil {
		logger.Error("error while Telegram bot initialization", log.Err(err))
		os.Exit(1)
	}

	logger.Info("successfully authorized in telegram api")

	removeWebhook := mustSetWebhook(ctx, logger, b, cfg.Bot.WebHookURL)
	defer removeWebhook()

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: b.WebhookHandler(),
	}

	go func() {
		logger.Info("starting http server", slog.String("address", server.Addr))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error while running http server", log.Err(err))
		}

		logger.Info("http server stopped")
	}()

	go b.StartWebhook(ctx)

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("error while HTTP server shutdown", log.Err(err))
	}
}

func parseAndValidateCLIArgs() (*CLIArgs, error) {
	configPath := flag.String("config", "", "path to config file")

	flag.Parse()

	// check if config path is not empty and file exists
	if *configPath == "" {
		return nil, errors.New("config path shouldn't be empty")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exists: %s", *configPath)
	}

	return &CLIArgs{configPath: *configPath}, nil
}
