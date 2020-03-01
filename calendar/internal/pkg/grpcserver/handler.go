package grpcserver

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"

	"github.com/golang/protobuf/ptypes"
)

// handler is gRPC server handler
type handler struct {
	logger  *zap.Logger
	storage domain.Storage
}

// Create new event
func (h *handler) Create(ctx context.Context, req *events.CreateRequest) (*events.Response, error) {
	event, err := protoToEvent(req.Event)
	if err != nil {
		h.logger.Error("parse event error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	err = h.storage.Insert(*event)
	if err != nil {
		h.logger.Error("insert to storage error", zap.Error(err))
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &events.Response{Error: events.ErrorCode_OK}, nil
}

// Update current event
func (h *handler) Update(ctx context.Context, req *events.UpdateRequest) (*events.Response, error) {
	event, err := protoToEvent(req.Event)
	if err != nil {
		h.logger.Error("parse event error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	err = h.storage.Update(*event)
	if err != nil {
		h.logger.Error("update to storage error", zap.Error(err))
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &events.Response{Error: events.ErrorCode_OK}, nil
}

// Remove event
func (h *handler) Remove(ctx context.Context, req *events.RemoveRequest) (*events.Response, error) {
	id, err := uuid.Parse(req.GetUuid())
	if err != nil {
		h.logger.Error("parse request error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	err = h.storage.Remove(domain.EventID(id))
	if err != nil {
		h.logger.Error("remove from storage error", zap.Error(err))
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &events.Response{Error: events.ErrorCode_OK}, nil
}

// DailyEventList get daily events
func (h *handler) DailyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

// WeeklyEventList get weekly events
func (h *handler) WeeklyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

// MonthlyEventList get monthly events
func (h *handler) MonthlyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

func protoToEvent(protoEvent *events.Event) (*domain.Event, error) {
	event := &domain.Event{
		Heading:     protoEvent.GetHeading(),
		Description: protoEvent.GetDescription(),
		Owner:       protoEvent.GetOwner(),
	}

	id, err := uuid.Parse(protoEvent.GetUuid())
	if err != nil {
		return nil, err
	}

	event.ID = domain.EventID(id)

	dateTime, err := ptypes.Timestamp(protoEvent.GetDateTime())
	if err != nil {
		return nil, err
	}

	event.DateTime = dateTime

	duration, err := ptypes.Duration(protoEvent.GetDuration())
	if err != nil {
		return nil, err
	}

	event.Duration = duration

	return event, nil
}

func eventToProto(protoEvent domain.Event) events.Event {
	return events.Event{}
}
