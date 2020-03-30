package scheduler

import (
	"context"
	"time"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"go.uber.org/zap"
)

// Scheduler is main struct
type Scheduler struct {
	logger    *zap.Logger
	storage   domain.Storage
	publisher domain.Producer
}

// New create new scheduler
func New(logger *zap.Logger, storage domain.Storage, publisher domain.Producer) *Scheduler {
	return &Scheduler{
		logger:    logger,
		storage:   storage,
		publisher: publisher,
	}
}

// Schedule and publish events
func (s *Scheduler) Schedule(ctx context.Context, duration time.Duration) error {
	last := time.Unix(0, 0)
	timeChan := time.Tick(duration)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timeChan:
			sDuration := time.Now().Sub(last)
			events, err := s.storage.GetEventsInTime(ctx, last, sDuration)
			if err != nil {
				s.logger.Error("storage GetEventsInTime error", zap.Error(err))
				return err
			}

			s.logger.Debug("publish", zap.Int("events length", len(events)))

			for _, event := range events {
				if err := s.publisher.Publish(event); err != nil {
					s.logger.Error("publish  error", zap.Error(err))
					return err
				}
			}

			last = last.Add(sDuration)
		}
	}
}
