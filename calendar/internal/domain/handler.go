package domain

// Handler handle events from delivery
type Handler func(data *Event) (string, error)
