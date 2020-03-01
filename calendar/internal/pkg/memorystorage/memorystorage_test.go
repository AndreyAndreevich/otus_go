package memorystorage

import (
	"testing"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_InsertDuplicate(t *testing.T) {
	storage := New()

	event := domain.Event{ID: domain.EventID(uuid.New())}

	err := storage.Insert(event)
	assert.NoError(t, err)

	err = storage.Insert(event)
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
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	events, err := storage.Listing()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestMemoryStorage_InsertWithZeroId(t *testing.T) {
	storage := New()
	storage.Insert(domain.Event{})

	err := storage.Insert(domain.Event{})
	assert.NoError(t, err)

	events, err := storage.Listing()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.NotEqual(t, events[0].ID, events[1].ID)
}

func TestMemoryStorage_RemoveEmptyStorage(t *testing.T) {
	storage := New()
	err := storage.Remove(domain.Event{ID: domain.EventID(uuid.New())})
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_RemoveIncorrectId(t *testing.T) {
	storage := New()
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	err := storage.Remove(domain.Event{ID: domain.EventID(uuid.New())})
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_Remove(t *testing.T) {
	storage := New()

	event := domain.Event{ID: domain.EventID(uuid.New())}

	storage.Insert(event)
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	err := storage.Remove(event)
	assert.NoError(t, err)

	events, _ := storage.Listing()
	assert.Len(t, events, 1)
}

func TestMemoryStorage_UpdateIncorrectId(t *testing.T) {
	storage := New()
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	err := storage.Update(domain.Event{ID: domain.EventID(uuid.New())})
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_Update(t *testing.T) {
	storage := New()

	event := domain.Event{ID: domain.EventID(uuid.New())}

	storage.Insert(event)

	err := storage.Update(domain.Event{ID: event.ID, Description: "new_data"})
	assert.NoError(t, err)

	events, _ := storage.Listing()
	assert.Len(t, events, 1)
	assert.Equal(t, "new_data", events[0].Description)
}
