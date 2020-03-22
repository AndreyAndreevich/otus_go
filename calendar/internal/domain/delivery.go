package domain

import "golang.org/x/net/context"

//go:generate mockery -name Delivery -output ../mocks

// Delivery is interface of events delivery
type Delivery interface {
	AddHandler(pattern string, handler Handler)
	// blocked method
	Run(ctx context.Context) error
}
