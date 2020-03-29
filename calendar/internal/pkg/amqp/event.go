package amqp

import (
	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
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

func eventToJson(event domain.Event) Event {
	return Event{
		ID:          event.ID.String(),
		Heading:     event.Heading,
		DateTime:    event.DateTime.Unix(),
		Duration:    int64(event.Duration.Seconds()),
		Description: event.Description,
		Owner:       event.Owner,
	}
}

func eventFromJson(event Event) (domain.Event, error) {
	return domain.Event{}, nil
}
