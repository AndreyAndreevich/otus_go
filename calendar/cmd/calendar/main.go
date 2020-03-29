package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/postgresstorage"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/config"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/grpcserver"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/httpserver"

	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/calendar"
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

	eventsDelivery := httpserver.New(logger, cfg.HTTPListen.IP, cfg.HTTPListen.Port)
	currentCalendar := calendar.New(logger, storage, eventsDelivery)
	gRPCServer := grpcserver.New(logger, cfg.GRPC.IP, cfg.GRPC.Port, currentCalendar)

	ctx, cancel := context.WithCancel(context.Background())

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := currentCalendar.Run(ctx); err != nil {
			logger.Error("error calendar run", zap.Error(err))
		}
		logger.Info("Calendar stopped")
		cancel()
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := gRPCServer.Run(ctx)
		if err != nil {
			logger.Error("gRPC server run error", zap.Error(err))
		}
		logger.Info("gRPC stopped")
		cancel()
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

	cfg.OutputPaths = []string{
		logFile,
	}

	cfg.Level = atom

	return cfg.Build()
}
