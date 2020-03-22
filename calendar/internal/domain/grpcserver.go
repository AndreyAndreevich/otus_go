package domain

import "context"

//go:generate mockery -name GRPCServer -output ../mocks

// GRPCServer is interface of gRPC server
type GRPCServer interface {
	Run(ctx context.Context) error
}
