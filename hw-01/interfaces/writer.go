package interfaces

type Writer interface {
	Write(args ... interface{}) error
}
