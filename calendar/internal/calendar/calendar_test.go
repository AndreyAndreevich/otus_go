package calendar

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/mocks"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testLayout = "2006-01-02"

var errTest = errors.New("test error")

func createCalendar() (*Calendar, *mocks.Storage) {
	storage := &mocks.Storage{}
	return &Calendar{
		logger:  zap.NewNop(),
		storage: storage,
	}, storage
}

func TestCalendar_CreateError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Insert", mock.Anything, mock.Anything).Return(errTest)

	err := calendar.Create(context.Background(), domain.Event{})

	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestCalendar_Create(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Insert", mock.Anything, mock.Anything).Return(nil)

	err := calendar.Create(context.Background(), domain.Event{})

	assert.NoError(t, err)
}

func TestCalendar_UpdateError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Update", mock.Anything, mock.Anything).Return(errTest)

	err := calendar.Update(context.Background(), domain.Event{})

	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestCalendar_Update(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := calendar.Update(context.Background(), domain.Event{})

	assert.NoError(t, err)
}

func TestCalendar_RemoveError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Remove", mock.Anything, mock.Anything).Return(errTest)

	err := calendar.Remove(context.Background(), domain.EventID{})

	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestCalendar_Remove(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("Remove", mock.Anything, mock.Anything).Return(nil)

	err := calendar.Remove(context.Background(), domain.EventID{})

	assert.NoError(t, err)
}

func TestCalendar_DailyEventListError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("GetEventsInTime", mock.Anything, mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := calendar.DailyEventList(context.Background(), time.Time{})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestCalendar_DailyEventList(t *testing.T) {
	calendar, storage := createCalendar()

	day, _ := time.Parse(testLayout, "2020-02-12")

	storage.On("GetEventsInTime", mock.Anything, day, time.Duration(time.Hour*24)).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := calendar.DailyEventList(context.Background(), day)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestCalendar_WeeklyEventListError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("GetEventsInTime", mock.Anything, mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := calendar.DailyEventList(context.Background(), time.Time{})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestHandler_WeeklyEventList(t *testing.T) {
	calendar, storage := createCalendar()

	day, _ := time.Parse(testLayout, "2020-02-12")

	storage.On("GetEventsInTime", mock.Anything, day, time.Duration(time.Hour*24*7)).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := calendar.WeeklyEventList(context.Background(), day)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestCalendar_MonthlyEventListError(t *testing.T) {
	calendar, storage := createCalendar()

	storage.On("GetEventsInTime", mock.Anything, mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := calendar.MonthlyEventList(context.Background(), time.Time{})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, errTest, err)
}

func TestCalendar_MonthlyEventList(t *testing.T) {
	calendar, storage := createCalendar()

	day, _ := time.Parse(testLayout, "2020-01-12")

	storage.On("GetEventsInTime", mock.Anything, day, time.Duration(time.Hour*24*31)).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := calendar.MonthlyEventList(context.Background(), day)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
