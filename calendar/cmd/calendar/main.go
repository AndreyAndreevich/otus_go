package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/sqlxstorage"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

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

	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)

	var cfg config.Config
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := newLogger(cfg.LogLvl, cfg.LogFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open("postgres", cfg.DB.DSN)
	if err != nil {
		logger.Error("connect to db error", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.DB.MaxConnections)
	db.SetMaxIdleConns(cfg.DB.MaxConnections)

	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Fatal("ping to db error", zap.Error(err))
		return
	}

	storage := sqlxstorage.New(logger, db)
	if err := storage.Migrate("postgres"); err != nil {
		logger.Fatal("migrate error", zap.Error(err))
	}

	eventsDelivery := httpserver.New(logger, cfg.HTTPListen.IP, cfg.HTTPListen.Port)
	gRPCServer := grpcserver.New(logger, cfg.GRPC.IP, cfg.GRPC.Port, storage)
	currentCalendar := calendar.New(logger, storage, eventsDelivery, gRPCServer)

	if err := currentCalendar.Run(); err != nil {
		logger.Fatal("error calendar run", zap.Error(err))
	}
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
