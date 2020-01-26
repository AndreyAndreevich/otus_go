package memorystorage

import (
	"testing"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/simpleevent"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/stretchr/testify/assert"
)

type TestUnknownEvent struct{}

func (e *TestUnknownEvent) GetID() domain.EventID {
	return domain.EventID(0)
}

func (e *TestUnknownEvent) SetID(id domain.EventID) {}

func (e *TestUnknownEvent) GetData() domain.EventData {
	return domain.EventData("data")
}

func TestMemoryStorage_InsertUnknownEvent(t *testing.T) {
	storage := New()
	err := storage.Insert(&TestUnknownEvent{})
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownEvent, err)
}

func TestMemoryStorage_UpdateUnknownEvent(t *testing.T) {
	storage := New()
	err := storage.Insert(&TestUnknownEvent{})
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownEvent, err)
}

func TestMemoryStorage_RemoveUnknownEvent(t *testing.T) {
	storage := New()
	err := storage.Insert(&TestUnknownEvent{})
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownEvent, err)
}

func TestMemoryStorage_InsertDuplicate(t *testing.T) {
	storage := New()
	err := storage.Insert(simpleevent.New(1, "data"))
	assert.NoError(t, err)

	err = storage.Insert(simpleevent.New(1, "data"))
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateEventID, err)
}

func TestMemoryStorage_ListingEmpty(t *testing.T) {
	storage := New()
	events, err := storage.Listing()
	assert.NoError(t, err)
	assert.Empty(t, events)
}

func TestMemoryStorage_Listing(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(1, "data1"))
	storage.Insert(simpleevent.New(2, "data2"))

	events, err := storage.Listing()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestMemoryStorage_InsertWithZeroId(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(0, "data"))

	err := storage.Insert(simpleevent.New(0, "data"))
	assert.NoError(t, err)

	events, err := storage.Listing()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.NotEqual(t, events[0].GetID(), events[1].GetID())
}

func TestMemoryStorage_RemoveEmptyStorage(t *testing.T) {
	storage := New()
	err := storage.Remove(simpleevent.New(1, "data"))
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_RemoveIncorrectId(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(1, "data"))
	storage.Insert(simpleevent.New(2, "data"))

	err := storage.Remove(simpleevent.New(3, "data"))
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_Remove(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(1, "data"))
	storage.Insert(simpleevent.New(2, "data"))

	err := storage.Remove(simpleevent.New(1, "data"))
	assert.NoError(t, err)

	events, _ := storage.Listing()
	assert.Len(t, events, 1)
	assert.Equal(t, domain.EventID(2), events[0].GetID())
}

func TestMemoryStorage_UpdateIncorrectId(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(1, "data"))

	err := storage.Update(simpleevent.New(2, "new_data"))
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_Update(t *testing.T) {
	storage := New()
	storage.Insert(simpleevent.New(1, "data"))

	err := storage.Update(simpleevent.New(1, "new_data"))
	assert.NoError(t, err)

	events, _ := storage.Listing()
	assert.Len(t, events, 1)
	assert.Equal(t, domain.EventData("new_data"), events[0].GetData())
}
