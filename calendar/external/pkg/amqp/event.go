package amqp

import (
	"time"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/google/uuid"
)

// Event - json event
type Event struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	DateTime    int64  `json:"date_time"`
	Duration    int64  `json:"duration_s"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
}

func eventToJSON(event domain.Event) Event {
	return Event{
		ID:          event.ID.String(),
		Heading:     event.Heading,
		DateTime:    event.DateTime.Unix(),
		Duration:    int64(event.Duration.Seconds()),
		Description: event.Description,
		Owner:       event.Owner,
	}
}

func eventFromJSON(event Event) (domain.Event, error) {
	id, err := uuid.Parse(event.ID)
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		ID:          id,
		Heading:     event.Heading,
		DateTime:    time.Unix(event.DateTime, 0),
		Duration:    time.Duration(event.Duration) * time.Second,
		Description: event.Description,
		Owner:       event.Owner,
	}, nil
}
