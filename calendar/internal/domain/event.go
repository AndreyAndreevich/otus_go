package domain

// EventID is id of events
type EventID uint64

// EventData is data from events
type EventData string

// Event is base interface of events
type Event interface {
	GetID() EventID
	SetID(id EventID)
	GetData() EventData
}
