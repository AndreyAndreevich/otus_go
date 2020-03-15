package calendar

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
)

// Calendar is main struct in program
type Calendar struct {
	logger     *zap.Logger
	storage    domain.Storage
	delivery   domain.Delivery
	gRPCServer domain.GRPCServer
}

// New create new calendar
func New(logger *zap.Logger, storage domain.Storage, delivery domain.Delivery, gRPCServer domain.GRPCServer) *Calendar {
	return &Calendar{
		logger:     logger,
		storage:    storage,
		delivery:   delivery,
		gRPCServer: gRPCServer,
	}
}

// Run calendar logic
func (c *Calendar) Run() error {
	c.logger.Info("I'm the best of calendars")

	c.delivery.AddHandler("/hello", func(data *domain.Event) (string, error) {
		c.logger.Info("request", zap.Reflect("data", data))
		return "world", nil
	})

	waitGroup := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := c.delivery.Run(ctx)
		if err != nil {
			c.logger.Error("delivery run error", zap.Error(err))
			cancel()
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := c.gRPCServer.Run(ctx)
		if err != nil {
			c.logger.Error("gRPC server run error", zap.Error(err))
			cancel()
		}
	}()

	waitGroup.Wait()

	return nil
}
