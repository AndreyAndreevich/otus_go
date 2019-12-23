package hw04

// Item is node of List
type Item struct {
	prev, next *Item
	v          interface{}
}

// Value returned value from item
func (i Item) Value() interface{} {
	return i.v
}

// Prev returned prev Item
func (i Item) Prev() *Item {
	return i.prev
}

// Next returned next Item
func (i Item) Next() *Item {
	return i.next
}
