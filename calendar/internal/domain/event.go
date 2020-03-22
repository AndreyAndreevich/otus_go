package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventID is id of events
type EventID = uuid.UUID

// Event is base interface of events
type Event struct {
	ID          EventID
	Heading     string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	Owner       string
}
