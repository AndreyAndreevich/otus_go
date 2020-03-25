package domain

import (
	"context"
	"time"
)

//go:generate mockery -name Calendar -output ../mocks

type Calendar interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, event Event) error
	Remove(ctx context.Context, id EventID) error
	DailyEventList(context.Context, time.Time) ([]Event, error)
	WeeklyEventList(context.Context, time.Time) ([]Event, error)
	MonthlyEventList(context.Context, time.Time) ([]Event, error)
}
