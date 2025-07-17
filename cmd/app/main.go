package main

import (
	"context"
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/config"
	eventconsummer "drillCore/internal/event-consummer"
	"drillCore/internal/events/telegram"
	"drillCore/internal/storage/debt/postgres"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err := setUpLogger(cfg.AppEnvs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	storage, err := postgres.New(ctx, cfg.DbEnvs, logger)

	eventsProcessor := telegram.New(
		tgClient.New(cfg.TelegramEnvs),
		storage,
	)

	logger.Info("Starting telegram bot")

	consumer := eventconsummer.New(eventsProcessor, eventsProcessor, cfg.TelegramEnvs.BatchSize)
	if err := consumer.Start(); err != nil {
		logger.Fatalf("service stopped:%v", err)
	}

}

func setUpLogger(cfg *config.AppEnvs) (*zap.SugaredLogger, error) {
	logConfig := zap.NewProductionConfig()

	logConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC1123))
	}

	switch cfg.Env {
	case "local":
		logConfig.Encoding = "console"
	case "dev":
		logConfig.Encoding = "console"
	case "prod":
		logConfig.Encoding = "json"
	}

	switch cfg.DebugFlag {
	case true:
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case false:
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
