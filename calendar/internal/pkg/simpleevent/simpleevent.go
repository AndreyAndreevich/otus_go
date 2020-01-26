package simpleevent

import "github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

// SimpleEvent - simple event realization
type SimpleEvent struct {
	id   domain.EventID
	data domain.EventData
}

// New - create new SimpleEvent
func New(id domain.EventID, data domain.EventData) *SimpleEvent {
	return &SimpleEvent{
		id:   id,
		data: data,
	}
}

// GetID - get event id
func (e *SimpleEvent) GetID() domain.EventID {
	return e.id
}

// SetID - set event id
func (e *SimpleEvent) SetID(id domain.EventID) {
	e.id = id
}

// GetData - get event data
func (e *SimpleEvent) GetData() domain.EventData {
	return e.data
}
