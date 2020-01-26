package domain

// Storage is base interface of events storage
type Storage interface {
	Insert(event Event) error
	Remove(event Event) error
	Update(event Event) error
	Listing() ([]Event, error)
}
