package mock

import (
	"time"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/stretchr/testify/mock"
)

// StorageMock is storage mock
type StorageMock struct {
	mock.Mock
}

// Insert as storage
func (s *StorageMock) Insert(event domain.Event) error {
	args := s.Called(event)
	return args.Error(0)
}

// Remove as storage
func (s *StorageMock) Remove(id domain.EventID) error {
	args := s.Called(id)
	return args.Error(0)
}

// Update as storage
func (s *StorageMock) Update(event domain.Event) error {
	args := s.Called(event)
	return args.Error(0)
}

// Listing as storage
func (s *StorageMock) Listing() ([]domain.Event, error) {
	args := s.Called()
	return args.Get(0).([]domain.Event), args.Error(1)
}

// GetEventsInTime as storage
func (s *StorageMock) GetEventsInTime(time time.Time, duration time.Duration) ([]domain.Event, error) {
	args := s.Called(time, duration)
	return args.Get(0).([]domain.Event), args.Error(1)
}
