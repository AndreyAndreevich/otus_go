package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/AndreyAndreevich/otus_go/calendar/external/pkg/amqp"

	_ "github.com/lib/pq"

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

	errorChan := make(chan error)
	waitGroup := &sync.WaitGroup{}

	consumer, err := amqp.NewConsumer(logger,
		errorChan,
		cfg.RabbitConfig.DSN,
		cfg.RabbitConfig.Exchange,
		cfg.RabbitConfig.Queue,
		waitGroup,
	)
	if err != nil {
		logger.Fatal("consumer create error", zap.Error(err))
	}

	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		select {
		case err, ok := <-errorChan:
			if ok {
				logger.Error("error from error channel", zap.Error(err))
				cancel()
			}
		case <-ctx.Done():

		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := consumer.Consume(ctx, func(event domain.Event) error {
			logger.Info("handle event", zap.Reflect("event", event))
			return nil
		})
		if err != nil {
			logger.Error("consume error", zap.Error(err))
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

	cfg.OutputPaths = []string{"stdout"}

	cfg.Level = atom

	return cfg.Build()
}
