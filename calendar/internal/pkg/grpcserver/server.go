package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/AndreyAndreevich/otus_go/calendar/external/pkg/events"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

// GRPCServer is gRPC server
type GRPCServer struct {
	logger   *zap.Logger
	addr     string
	calendar domain.Calendar
}

// New created new GRPCServer
func New(logger *zap.Logger, ip string, port int, calendar domain.Calendar) *GRPCServer {
	return &GRPCServer{
		logger:   logger,
		addr:     fmt.Sprintf("%s:%d", ip, port),
		calendar: calendar,
	}
}

// Run GRPCServer
func (s *GRPCServer) Run(ctx context.Context) error {
	s.logger.Debug("gRPC server starting")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error("Listening error", zap.Error(err))
		return err
	}

	server := grpc.NewServer()
	reflection.Register(server) // for evans
	events.RegisterGRPCServer(server, &handler{
		logger:  s.logger,
		storage: s.storage,
	})

	go func(ctx context.Context, srv *grpc.Server) {
		<-ctx.Done()
		server.GracefulStop()
	}(ctx, server)

	return server.Serve(lis)
}
