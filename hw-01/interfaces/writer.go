package interfaces

// Writer is interface for write info to somewhere
type Writer interface {
	Write(args ...interface{}) error
}
