package amqp

import "github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

type Event struct {
}

func eventToJson(event domain.Event) (Event, error) {
	return Event{}, nil
}

func eventFromJson(event Event) (domain.Event, error) {
	return domain.Event{}, nil
}
