package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"protomorphine/tg-notes/internal/bot/handlers/fallback"
	"protomorphine/tg-notes/internal/bot/handlers/help"
	"protomorphine/tg-notes/internal/bot/middleware"
	"protomorphine/tg-notes/internal/config"
	sl "protomorphine/tg-notes/internal/logger"

	"github.com/go-telegram/bot"
)

const (
	// log levels
	Debug string = "DEBUG"
	Info  string = "INFO"
	Warn  string = "WARN"
	Error string = "ERROR"

	// application enviroments
	Local      string = "local"
	Production string = "prod"
)

type CLIArgs struct {
	configPath string
}

func main() {
	args := mustParseCLIArgs()

	cfg, err := config.Load(args.configPath)
	if err != nil {
		slog.Error("error occured while loading config", slog.Any("err", err))
		os.Exit(1)
	}

	logger := configureLogger(cfg.Environment, cfg.Logger)
	logger = logger.With(slog.String("env", cfg.Environment))

	logger.Info("starting tg-notes app")

	opts := []bot.Option{
		bot.WithDefaultHandler(fallback.New(logger)),
		bot.WithCheckInitTimeout(cfg.Bot.InitTimeout),
		bot.WithMiddlewares(
			middleware.NewReqID(),
			middleware.NewLog(logger),
		),
	}

	b, err := bot.New(cfg.Bot.Key, opts...)
	if err != nil {
		logger.Error("unable to initialize bot", sl.Err(err))
		os.Exit(1)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, help.New(logger))
	b.Start(context.TODO())
}

func mustParseCLIArgs() *CLIArgs {
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

func configureLogger(env string, logCfg config.LoggerConfig) *slog.Logger {
	var level slog.Level

	switch logCfg.MinLevel {
	case Info:
		level = slog.LevelInfo
	case Debug:
		level = slog.LevelDebug
	case Warn:
		level = slog.LevelWarn
	case Error:
		level = slog.LevelError
	}

	handlerOptions := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch env {
	case Local:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	case Production:
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	return slog.New(handler)
}
