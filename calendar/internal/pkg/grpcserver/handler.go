package grpcserver

import (
	"context"
	"fmt"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"
)

type Handler struct {
}

func (s *Handler) Create(ctx context.Context, req *events.CreateRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

func (s *Handler) Update(ctx context.Context, req *events.UpdateRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

func (s *Handler) Remove(ctx context.Context, req *events.RemoveRequest) (*events.Response, error) {
	fmt.Println(req)
	return nil, nil
}

func (s *Handler) DailyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

func (s *Handler) WeeklyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}

func (s *Handler) MonthlyEventList(ctx context.Context, req *events.DataRequest) (*events.EventsResponse, error) {
	fmt.Println(req)
	return nil, nil
}
