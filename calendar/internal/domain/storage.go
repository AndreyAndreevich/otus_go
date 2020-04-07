package domain

import (
	"context"
	"time"
)

//go:generate mockery -name Storage -output ../mocks

// ScheduleStorage is interface of schedule events storage
type ScheduleStorage interface {
	GetEventsInTime(ctx context.Context, time time.Time, duration time.Duration) ([]Event, error)
}

// Storage is base interface of events storage
type Storage interface {
	ScheduleStorage
	Insert(ctx context.Context, event Event) error
	Remove(ctx context.Context, id EventID) error
	Update(ctx context.Context, event Event) error
	Listing(ctx context.Context) ([]Event, error)
	Close() error
}
