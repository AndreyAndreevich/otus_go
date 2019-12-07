package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type ClockMock struct {
	mock.Mock
}

func (c *ClockMock) GetCurrentTime() (time.Time, error) {
	args := c.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

type WriterMock struct {
	mock.Mock
	currentTime time.Time
}

func (w *WriterMock) Write(args ...interface{}) error {
	w.currentTime = args[0].(time.Time)
	return nil
}

func TestApp_Run_With_Error(t *testing.T) {
	clock := &ClockMock{}
	writer := &WriterMock{}
	app := NewApp(clock, writer)

	clock.On("GetCurrentTime").Return(time.Now(), errors.New("bad time"))

	err := app.Run()
	assert.Error(t, err)
}

func TestApp_Run_WriteTime(t *testing.T) {
	clock := &ClockMock{}
	writer := &WriterMock{}
	app := NewApp(clock, writer)

	currentTime := time.Now()
	clock.On("GetCurrentTime").Return(currentTime, nil)

	err := app.Run()
	assert.NoError(t, err)
	assert.Equal(t, currentTime, writer.currentTime)
}
