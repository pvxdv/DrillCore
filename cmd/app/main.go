package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"drillCore/internal/bot"
	"drillCore/internal/config"
	"drillCore/internal/events/event-consummer"
	"drillCore/internal/events/event-processor"
	"drillCore/internal/events/event-processor/manager"
	"drillCore/internal/events/event-processor/manager/command"
	"drillCore/internal/events/event-processor/manager/date"
	"drillCore/internal/events/event-processor/manager/debt"
	mainmenu "drillCore/internal/events/event-processor/manager/main-menu"
	"drillCore/internal/session"
	"drillCore/internal/storage/debt/postgres"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	logger.Debugf("resived config: %+v", cfg)

	storage, err := postgres.New(ctx, cfg.DbEnvs, logger)

	tg := bot.New(cfg.TelegramEnvs, logger)

	sMng := session.New()

	debtH := debt.New(tg, sMng, storage, logger)
	cmdH := command.New(tg, sMng, logger)
	menuH := mainmenu.New(tg, sMng, logger)
	dateH := date.New(tg, sMng, logger)

	hMng := manager.New(tg, sMng, logger, cmdH, menuH, debtH, dateH)

	eventsProcessor := eventprocessor.New(tg, hMng, logger)

	logger.Info("Starting event-processor bot")

	consumer := eventconsummer.New(eventsProcessor, eventsProcessor, cfg.TelegramEnvs.BatchSize, logger)
	if err := consumer.Start(ctx); err != nil {
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
