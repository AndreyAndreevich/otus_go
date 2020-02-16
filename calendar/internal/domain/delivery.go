package domain

// Delivery is interface of events delivery
type Delivery interface {
	AddHandler(pattern string, handler Handler)
	// blocked method
	Run() error
}
