package memorystorage

import (
	"errors"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/simpleevent"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
)

// EventType is type of events
type EventType int

const (
	// UnknownEvent - unknown event type
	UnknownEvent = EventType(0)
	// SimpleEventType - type of SimpleEvent
	SimpleEventType = EventType(1)
)

var (
	// ErrUnknownEvent - unknown event error
	ErrUnknownEvent = errors.New("unknown event")
	// ErrDuplicateEventID - duplicate event id error
	ErrDuplicateEventID = errors.New("duplicate event id")
	// ErrNotExist - not exist event error
	ErrNotExist = errors.New("event is not exist")
)

type eventData struct {
	eventType EventType
	eventData domain.EventData
}

// MemoryStorage - event's storage in memory
type MemoryStorage struct {
	data map[domain.EventID]eventData
}

// New create new MemoryStorage
func New() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[domain.EventID]eventData),
	}
}

// Insert event into MemoryStorage
func (s *MemoryStorage) Insert(event domain.Event) error {
	eventType := checkEventType(event)
	if eventType == UnknownEvent {
		return ErrUnknownEvent
	}

	if event.GetID() == 0 {
		nextID := domain.EventID(1)
		for id := range s.data {
			if id >= nextID {
				nextID = id + 1
			}
		}
		event.SetID(nextID)
	} else {
		_, isExist := s.data[event.GetID()]
		if isExist {
			return ErrDuplicateEventID
		}
	}

	s.data[event.GetID()] = eventData{eventType, event.GetData()}

	return nil
}

// Remove event from MemoryStorage
func (s *MemoryStorage) Remove(event domain.Event) error {
	eventType := checkEventType(event)
	if eventType == UnknownEvent {
		return ErrUnknownEvent
	}

	_, isExist := s.data[event.GetID()]
	if !isExist {
		return ErrNotExist
	}

	delete(s.data, event.GetID())

	return nil
}

// Update event in MemoryStorage
func (s *MemoryStorage) Update(event domain.Event) error {
	eventType := checkEventType(event)
	if eventType == UnknownEvent {
		return ErrUnknownEvent
	}

	data, isExist := s.data[event.GetID()]
	if !isExist {
		return ErrNotExist
	}

	data.eventType = eventType
	data.eventData = event.GetData()

	s.data[event.GetID()] = data

	return nil
}

// Listing - get all events from MemoryStorage
func (s *MemoryStorage) Listing() ([]domain.Event, error) {
	events := make([]domain.Event, 0, len(s.data))

	for id, data := range s.data {
		var event domain.Event

		switch data.eventType {
		case SimpleEventType:
			event = simpleevent.New(id, data.eventData)
		default:
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func checkEventType(event domain.Event) EventType {
	switch event.(type) {
	case *simpleevent.SimpleEvent:
		return SimpleEventType
	default:
		return UnknownEvent
	}
}
