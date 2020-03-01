package domain

// Storage is base interface of events storage
type Storage interface {
	Insert(event Event) error
	Remove(id EventID) error
	Update(event Event) error
	Listing() ([]Event, error)
}
