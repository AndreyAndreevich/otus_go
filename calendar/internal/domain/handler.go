package domain

// Handler handle events from delivery
type Handler func(data EventData) (string, error)
