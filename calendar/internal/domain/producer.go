package domain

//go:generate mockery -name Producer -output ../mocks

// Producer is interface for async publish events
type Producer interface {
	Publish(event Event) error
	Close() error
}
