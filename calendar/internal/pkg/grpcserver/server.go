package grpcserver

import (
	"fmt"
	"net"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

// GRPCServer is gRPC server
type GRPCServer struct {
	logger  *zap.Logger
	addr    string
	storage domain.Storage
}

// New created new GRPCServer
func New(logger *zap.Logger, ip string, port int, storage domain.Storage) *GRPCServer {
	return &GRPCServer{
		logger:  logger,
		addr:    fmt.Sprintf("%s:%d", ip, port),
		storage: storage,
	}
}

// Run GRPCServer
func (s *GRPCServer) Run() error {
	s.logger.Debug("gRPC server starting")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error("Listening error", zap.Error(err))
		return err
	}

	server := grpc.NewServer()
	events.RegisterGRPCServer(server, &handler{
		logger:  s.logger,
		storage: s.storage,
	})
	return server.Serve(lis)
}