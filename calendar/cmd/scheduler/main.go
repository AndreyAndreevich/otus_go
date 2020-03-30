package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/scheduler"

	"github.com/AndreyAndreevich/otus_go/calendar/external/pkg/amqp"

	_ "github.com/lib/pq"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/postgresstorage"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/config"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "", "path to config")

	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := newLogger(cfg.LogLvl, cfg.LogFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := postgresstorage.New(logger, cfg.DB.DSN, cfg.DB.MaxConnections, cfg.DB.IdleConnections)
	if err != nil {
		logger.Fatal("db create error", zap.Error(err))
	}
	defer storage.Close()

	if err := storage.HealthCheck(); err != nil {
		logger.Fatal("HealthCheck storage error", zap.Error(err))
	}

	errorChan := make(chan error)
	waitGroup := sync.WaitGroup{}

	publisher, err := amqp.NewProducer(logger, errorChan, cfg.RabbitConfig.DSN, cfg.RabbitConfig.Exchange, &waitGroup)
	if err != nil {
		logger.Fatal("create publisher error", zap.Error(err))
	}
	defer publisher.Close()

	currentScheduler := scheduler.New(logger, storage, publisher)

	ctx, cancel := context.WithCancel(context.Background())

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := currentScheduler.Schedule(ctx, cfg.ScheduleDuration); err != nil {
			logger.Error("error schedule", zap.Error(err))
		}
		logger.Info("Scheduler stopped")
		cancel()
	}()

	waitGroup.Add(1)
	go func() {
		select {
		case <-ctx.Done():
		case err, ok := <-errorChan:
			if ok {
				logger.Error("error from error channel", zap.Error(err))
				cancel()
			}
		}
	}()

	waitGroup.Wait()
}

func newLogger(level, logFile string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	atom := zap.NewAtomicLevel()
	err := atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	if level != "debug" {
		cfg.OutputPaths = []string{logFile}
	} else {
		cfg.OutputPaths = []string{"stdout"}
	}

	cfg.Level = atom

	return cfg.Build()
}
