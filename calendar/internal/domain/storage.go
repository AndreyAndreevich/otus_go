package domain

import (
	"context"
	"time"
)

//go:generate mockery -name Storage -output ../mocks

// Storage is base interface of events storage
type Storage interface {
	Insert(ctx context.Context, event Event) error
	Remove(ctx context.Context, id EventID) error
	Update(ctx context.Context, event Event) error
	Listing(ctx context.Context) ([]Event, error)
	GetEventsInTime(ctx context.Context, time time.Time, duration time.Duration) ([]Event, error)
}
