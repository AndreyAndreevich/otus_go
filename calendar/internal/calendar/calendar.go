package calendar

import (
	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
)

// Calendar is main struct in program
type Calendar struct {
	logger  *zap.Logger
	storage domain.Storage
}

// New create new calendar
func New(logger *zap.Logger, storage domain.Storage) *Calendar {
	return &Calendar{
		logger:  logger,
		storage: storage,
	}
}

// Run calendar logic
func (c *Calendar) Run() error {
	c.logger.Info("I'm the best of calendars")

	return nil
}
