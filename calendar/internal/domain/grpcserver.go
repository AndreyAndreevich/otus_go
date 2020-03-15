package domain

import "context"

// GRPCServer is interface of gRPC server
type GRPCServer interface {
	Run(ctx context.Context) error
}
