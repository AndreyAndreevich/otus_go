package grpcserver

import (
	"fmt"
	"net"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/events"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

type GRPCServer struct {
	logger *zap.Logger
	addr   string
}

func New(logger *zap.Logger, ip string, port int) *GRPCServer {
	return &GRPCServer{
		logger: logger,
		addr:   fmt.Sprintf("%s:%d", ip, port),
	}
}

func (s *GRPCServer) Run() error {
	s.logger.Debug("gRPC server starting")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error("Listening error", zap.Error(err))
		return err
	}

	server := grpc.NewServer()
	events.RegisterGRPCServer(server, &Handler{})
	return server.Serve(lis)
}
