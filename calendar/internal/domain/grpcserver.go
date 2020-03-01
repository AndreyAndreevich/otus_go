package domain

// GRPCServer is interface of gRPC server
type GRPCServer interface {
	Run() error
}
