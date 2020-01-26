package calendar

import (
	"fmt"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
)

// Calendar is main struct in program
type Calendar struct {
	storage domain.Storage
}

// New create new calendar
func New(storage domain.Storage) *Calendar {
	return &Calendar{
		storage: storage,
	}
}

// Run calendar logic
func (c *Calendar) Run() error {
	fmt.Println("I'm the best of calendars")

	return nil
}
