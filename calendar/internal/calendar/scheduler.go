package calendar

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Schedule publish events
func (c *Calendar) Schedule(ctx context.Context, duration time.Duration) error {
	last := time.Unix(0, 0)
	timeChan := time.Tick(duration)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timeChan:
			sDuration := time.Now().Sub(last)
			events, err := c.storage.GetEventsInTime(ctx, last, sDuration)
			if err != nil {
				c.logger.Error("storage GetEventsInTime error", zap.Error(err))
				return err
			}

			c.logger.Debug("publish", zap.Int("events length", len(events)))

			for _, event := range events {
				if err := c.publisher.Publish(event); err != nil {
					c.logger.Error("publish  error", zap.Error(err))
					return err
				}
			}

			last = last.Add(sDuration)
		}
	}
}
