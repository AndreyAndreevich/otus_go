package hw04

// List is linked list structure
type List struct {
	first, last *Item
	len         int
}

// Len returned length of list
func (l List) Len() int {
	return l.len
}

// First returned first Item
func (l List) First() *Item {
	return l.first
}

// Last returned last Item
func (l List) Last() *Item {
	return l.last
}

func (l *List) insertAfter(curItem, newItem *Item) {
	newItem.prev = curItem
	if curItem.next == nil {
		newItem.next = nil
		l.last = newItem
	} else {
		newItem.next = curItem.next
		curItem.next.prev = newItem
	}
	curItem.next = newItem
}

func (l *List) insertBefore(curItem, newItem *Item) {
	newItem.next = curItem
	if curItem.prev == nil {
		newItem.prev = nil
		l.first = newItem
	} else {
		newItem.prev = curItem.prev
		curItem.prev.next = newItem
	}
	curItem.prev = newItem
}

// PushFront added value to front of List
func (l *List) PushFront(v interface{}) {
	newItem := &Item{nil, l.first, v}

	if l.first == nil {
		l.first = newItem
		l.last = newItem
		newItem.prev = nil
		newItem.next = nil
	} else {
		l.insertBefore(l.first, newItem)
	}

	l.len++
}

// PushBack added value to back of List
func (l *List) PushBack(v interface{}) {
	if l.last == nil {
		l.PushFront(v)
	} else {
		newItem := &Item{l.last, nil, v}
		l.insertAfter(l.last, newItem)
		l.len++
	}
}

// Remove removed Item from List
// Item must be had in the List
func (l *List) Remove(item *Item) {
	if item.prev == nil {
		l.first = item.next
	} else {
		item.prev.next = item.next
	}

	if item.next == nil {
		l.last = item.prev
	} else {
		item.next.prev = item.prev
	}
	l.len--
}
