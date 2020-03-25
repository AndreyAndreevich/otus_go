package calendar

import (
	"context"
	"sync"
	"time"

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

func (c *Calendar) Create(ctx context.Context, event domain.Event) error {
	err := c.storage.Insert(ctx, event)
	if err != nil {
		c.logger.Error("insert to storage error", zap.Error(err))
		return err
	}
	return nil
}

func (c *Calendar) Update(ctx context.Context, event domain.Event) error {
	err := c.storage.Update(ctx, event)
	if err != nil {
		c.logger.Error("update to storage error", zap.Error(err))
		return err
	}
	return nil
}

func (c *Calendar) Remove(ctx context.Context, id domain.EventID) error {
	err := c.storage.Remove(ctx, id)
	if err != nil {
		c.logger.Error("remove from storage error", zap.Error(err))
		return err
	}
	return nil
}

func (c *Calendar) DailyEventList(ctx context.Context, dateTime time.Time) ([]domain.Event, error) {
	duration := time.Duration(time.Hour * 24)
	events, err := c.storage.GetEventsInTime(ctx, dateTime, duration)
	if err != nil {
		c.logger.Error("get events in time from storage error", zap.Error(err))
		return nil, err
	}
	return events, nil
}

func (c *Calendar) WeeklyEventList(ctx context.Context, dateTime time.Time) ([]domain.Event, error) {
	duration := time.Duration(time.Hour * 24 * 7)
	events, err := c.storage.GetEventsInTime(ctx, dateTime, duration)
	if err != nil {
		c.logger.Error("get events in time from storage error", zap.Error(err))
		return nil, err
	}
	return events, nil
}

func (c *Calendar) MonthlyEventList(ctx context.Context, dateTime time.Time) ([]domain.Event, error) {
	currentYear, currentMonth, _ := dateTime.Date()
	currentLocation := dateTime.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstOfMonthNextMonth := firstOfMonth.AddDate(0, 1, 0)

	duration := firstOfMonthNextMonth.Sub(firstOfMonth)

	events, err := c.storage.GetEventsInTime(ctx, dateTime, duration)
	if err != nil {
		c.logger.Error("get events in time from storage error", zap.Error(err))
		return nil, err
	}

	return events, nil
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
