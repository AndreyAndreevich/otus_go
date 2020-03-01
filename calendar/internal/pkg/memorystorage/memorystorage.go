package memorystorage

import (
	"errors"
	"sync"

	"github.com/google/uuid"

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

// MemoryStorage - event's storage in memory
type MemoryStorage struct {
	mtx  sync.Mutex
	data map[domain.EventID]domain.Event
}

// New create new MemoryStorage
func New() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[domain.EventID]domain.Event),
	}
}

// Insert event into MemoryStorage
func (s *MemoryStorage) Insert(event domain.Event) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if event.ID == domain.EventID(uuid.Nil) {
		event.ID = domain.EventID(uuid.New())
	} else {
		_, isExist := s.data[event.ID]
		if isExist {
			return ErrDuplicateEventID
		}
	}

	s.data[event.ID] = event

	return nil
}

// Remove event from MemoryStorage
func (s *MemoryStorage) Remove(id domain.EventID) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, isExist := s.data[id]
	if !isExist {
		return ErrNotExist
	}

	delete(s.data, id)

	return nil
}

// Update event in MemoryStorage
func (s *MemoryStorage) Update(event domain.Event) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, isExist := s.data[event.ID]
	if !isExist {
		return ErrNotExist
	}

	s.data[event.ID] = event

	return nil
}

// Listing - get all events from MemoryStorage
func (s *MemoryStorage) Listing() ([]domain.Event, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	events := make([]domain.Event, 0, len(s.data))

	for _, event := range s.data {
		events = append(events, event)
	}

	return events, nil
}
