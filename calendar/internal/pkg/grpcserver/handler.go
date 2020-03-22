package grpcserver

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

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

	err = h.storage.Insert(ctx, *event)
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

	err = h.storage.Update(ctx, *event)
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

	err = h.storage.Remove(ctx, domain.EventID(id))
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
	duration := time.Duration(time.Hour * 24)
	result, err := h.getEventsInTime(ctx, req.GetDateTime(), duration)
	if err != nil {
		return nil, err
	}

	return &events.EventsResponse{
		Error:  events.ErrorCode_OK,
		Events: result,
	}, nil
}

// WeeklyEventList get weekly events
func (h *handler) WeeklyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	duration := time.Duration(time.Hour * 24 * 7)
	result, err := h.getEventsInTime(ctx, req.GetDateTime(), duration)
	if err != nil {
		return nil, err
	}

	return &events.EventsResponse{
		Error:  events.ErrorCode_OK,
		Events: result,
	}, nil
}

// MonthlyEventList get monthly events
func (h *handler) MonthlyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	dateTime, err := ptypes.Timestamp(req.GetDateTime())
	if err != nil {
		h.logger.Error("parse timestamp error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	currentYear, currentMonth, _ := dateTime.Date()
	currentLocation := dateTime.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstOfMonthNextMonth := firstOfMonth.AddDate(0, 1, 0)

	duration := firstOfMonthNextMonth.Sub(firstOfMonth)

	result, err := h.getEventsInTime(ctx, req.GetDateTime(), duration)
	if err != nil {
		return nil, err
	}

	return &events.EventsResponse{
		Error:  events.ErrorCode_OK,
		Events: result,
	}, nil
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

func (h *handler) getEventsInTime(ctx context.Context,
	time *timestamp.Timestamp,
	duration time.Duration) ([]*events.Event, error) {

	dateTime, err := ptypes.Timestamp(time)
	if err != nil {
		h.logger.Error("parse timestamp error", zap.Error(err))
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	domEvents, err := h.storage.GetEventsInTime(ctx, dateTime, duration)
	if err != nil {
		h.logger.Error("get events in time from storage error", zap.Error(err))
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	result := []*events.Event{}
	for _, event := range domEvents {
		dateTime, err := ptypes.TimestampProto(event.DateTime)
		if err != nil {
			h.logger.Warn("incorrect datetime from event", zap.Reflect("event", event))
			continue
		}

		result = append(result, &events.Event{
			Uuid:        uuid.UUID(event.ID).String(),
			Heading:     event.Heading,
			DateTime:    dateTime,
			Duration:    ptypes.DurationProto(event.Duration),
			Description: event.Description,
			Owner:       event.Owner,
		})
	}

	return result, nil
}
