package grpcserver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"

	storagemock "github.com/AndreyAndreevich/otus_go/calendar/internal/mock"
)

const testLayout = "2006-01-02"

var errTest = errors.New("test error")

func createHandler() (*handler, *storagemock.StorageMock) {
	storage := &storagemock.StorageMock{}
	return &handler{
		logger:  zap.NewNop(),
		storage: storage,
	}, storage
}

func createEvent() *events.Event {
	return &events.Event{
		Uuid:     uuid.New().String(),
		DateTime: &timestamp.Timestamp{},
		Duration: &duration.Duration{},
	}
}

func TestHandler_CreateError(t *testing.T) {
	handler, storage := createHandler()

	req := &events.CreateRequest{Event: createEvent()}
	storage.On("Insert", mock.Anything).Return(errTest)

	res, err := handler.Create(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Create(t *testing.T) {
	handler, storage := createHandler()

	req := &events.CreateRequest{Event: createEvent()}
	storage.On("Insert", mock.Anything).Return(nil)

	res, err := handler.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_UpdateError(t *testing.T) {
	handler, storage := createHandler()

	req := &events.UpdateRequest{Event: createEvent()}
	storage.On("Update", mock.Anything).Return(errTest)

	res, err := handler.Update(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Update(t *testing.T) {
	handler, storage := createHandler()

	req := &events.UpdateRequest{Event: createEvent()}
	storage.On("Update", mock.Anything).Return(nil)

	res, err := handler.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_RemoveError(t *testing.T) {
	handler, storage := createHandler()

	req := &events.RemoveRequest{Uuid: uuid.New().String()}
	storage.On("Remove", mock.Anything).Return(errTest)

	res, err := handler.Remove(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_Remove(t *testing.T) {
	handler, storage := createHandler()

	req := &events.RemoveRequest{Uuid: uuid.New().String()}
	storage.On("Remove", mock.Anything).Return(nil)

	res, err := handler.Remove(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)
}

func TestHandler_DailyEventList(t *testing.T) {
	handler, storage := createHandler()

	day, _ := time.Parse(testLayout, "2020-02-12")

	req := &events.DataRequest{
		DateTime: &timestamp.Timestamp{Seconds: day.Unix()},
	}

	storage.On("GetEvents", mock.Anything).Return(nil)

	res, err := handler.Remove(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, events.ErrorCode_OK, res.Error)

}
