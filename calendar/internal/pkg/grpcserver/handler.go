package grpcserver

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/external/pkg/events"
	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/golang/protobuf/ptypes"
)

// handler is gRPC server handler
type handler struct {
	logger   *zap.Logger
	calendar domain.Calendar
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

	err = h.calendar.Create(ctx, *event)
	if err != nil {
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

	err = h.calendar.Update(ctx, *event)
	if err != nil {
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

	err = h.calendar.Remove(ctx, domain.EventID(id))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &events.Response{Error: events.ErrorCode_OK}, nil
}

// DailyEventList get daily events
func (h *handler) DailyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	return h.durationEventList(ctx, req, h.calendar.DailyEventList)
}

// WeeklyEventList get weekly events
func (h *handler) WeeklyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	return h.durationEventList(ctx, req, h.calendar.WeeklyEventList)
}

// MonthlyEventList get monthly events
func (h *handler) MonthlyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	return h.durationEventList(ctx, req, h.calendar.MonthlyEventList)
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

func (h *handler) eventToProto(event domain.Event) *events.Event {
	dateTime, err := ptypes.TimestampProto(event.DateTime)
	if err != nil {
		h.logger.Warn("incorrect datetime from event", zap.Reflect("event", event))
		return nil
	}

	return &events.Event{
		Uuid:        uuid.UUID(event.ID).String(),
		Heading:     event.Heading,
		DateTime:    dateTime,
		Duration:    ptypes.DurationProto(event.Duration),
		Description: event.Description,
		Owner:       event.Owner,
	}
}

type durationHandler = func(ctx context.Context, date time.Time) ([]domain.Event, error)

func (h *handler) durationEventList(ctx context.Context, req *events.DataRequest,
	dh durationHandler) (*events.EventsResponse, error) {

	dateTime, err := ptypes.Timestamp(req.GetDateTime())
	if err != nil {
		h.logger.Error("parse timestamp error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	domEvents, err := dh(ctx, dateTime)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	var protoEvents []*events.Event
	for _, event := range domEvents {
		protoEvent := h.eventToProto(event)
		if protoEvent != nil {
			protoEvents = append(protoEvents, protoEvent)
		}
	}

	return &events.EventsResponse{
		Error:  events.ErrorCode_OK,
		Events: protoEvents,
	}, nil
}
