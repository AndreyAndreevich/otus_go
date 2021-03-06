package grpcserver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/external/pkg/events"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/mocks"
)

const testLayout = "2006-01-02"

var errTest = errors.New("test error")

func createHandler() (*handler, *mocks.Calendar) {
	calendar := &mocks.Calendar{}
	return &handler{
		logger:   zap.NewNop(),
		calendar: calendar,
	}, calendar
}

func createEvent() *events.Event {
	return &events.Event{
		Uuid:     uuid.New().String(),
		DateTime: &timestamp.Timestamp{},
		Duration: &duration.Duration{},
	}
}

func TestHandler_CreateError(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.CreateRequest{Event: createEvent()}
	calendar.On("Create", mock.Anything, mock.Anything).Return(errTest)

	res, err := handler.Create(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Create(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.CreateRequest{Event: createEvent()}
	calendar.On("Create", mock.Anything, mock.Anything).Return(nil)

	res, err := handler.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_UpdateError(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.UpdateRequest{Event: createEvent()}
	calendar.On("Update", mock.Anything, mock.Anything).Return(errTest)

	res, err := handler.Update(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Update(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.UpdateRequest{Event: createEvent()}
	calendar.On("Update", mock.Anything, mock.Anything).Return(nil)

	res, err := handler.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_RemoveError(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.RemoveRequest{Uuid: uuid.New().String()}
	calendar.On("Remove", mock.Anything, mock.Anything).Return(errTest)

	res, err := handler.Remove(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Remove(t *testing.T) {
	handler, calendar := createHandler()

	req := &events.RemoveRequest{Uuid: uuid.New().String()}
	calendar.On("Remove", mock.Anything, mock.Anything).Return(nil)

	res, err := handler.Remove(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_DailyEventListError(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}
	calendar.On("DailyEventList", mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := handler.DailyEventList(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_DailyEventList(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}

	calendar.On("DailyEventList", mock.Anything, day).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := handler.DailyEventList(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, res.Events, 2)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_WeeklyEventListError(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}
	calendar.On("WeeklyEventList", mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := handler.WeeklyEventList(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_WeeklyEventList(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}

	calendar.On("WeeklyEventList", mock.Anything, day).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := handler.WeeklyEventList(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, res.Events, 2)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_MonthlyEventListError(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}
	calendar.On("MonthlyEventList", mock.Anything, mock.Anything).Return(nil, errTest)

	res, err := handler.MonthlyEventList(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_MonthlyEventList(t *testing.T) {
	handler, calendar := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")
	req := &events.DataRequest{DateTime: &timestamp.Timestamp{Seconds: day.Unix()}}

	calendar.On("MonthlyEventList", mock.Anything, day).Return([]domain.Event{
		{},
		{},
	}, nil)

	res, err := handler.MonthlyEventList(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, res.Events, 2)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}
