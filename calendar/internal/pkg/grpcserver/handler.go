package grpcserver

import (
	"context"
	"fmt"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"
)

// Handler is gRPC server handler
type Handler struct {
}

// Create new event
func (s *Handler) Create(ctx context.Context, req *events.CreateRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

// Update current event
func (s *Handler) Update(ctx context.Context, req *events.UpdateRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

// Remove event
func (s *Handler) Remove(ctx context.Context, req *events.RemoveRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

// DailyEventList get daily events
func (s *Handler) DailyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

// WeeklyEventList get weekly events
func (s *Handler) WeeklyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

// MonthlyEventList get monthly events
func (s *Handler) MonthlyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}
