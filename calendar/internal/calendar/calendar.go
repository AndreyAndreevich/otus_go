package calendar

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
)

// Calendar is main struct in program
type Calendar struct {
	logger   *zap.Logger
	storage  domain.Storage
	delivery domain.Delivery
}

// New create new calendar
func New(logger *zap.Logger, storage domain.Storage, delivery domain.Delivery) *Calendar {
	return &Calendar{
		logger:   logger,
		storage:  storage,
		delivery: delivery,
	}
}

// Run calendar logic
func (c *Calendar) Run(ctx context.Context) error {
	c.logger.Info("I'm the best of calendars")

	c.delivery.AddHandler("/hello", func(data *domain.Event) (string, error) {
		c.logger.Info("request", zap.Reflect("data", data))
		return "world", nil
	})

	waitGroup := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := c.delivery.Run(ctx)
		if err != nil {
			c.logger.Error("delivery run error", zap.Error(err))
			cancel()
		}
	}()

	waitGroup.Wait()

	return nil
}
