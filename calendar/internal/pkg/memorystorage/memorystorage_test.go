package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/stretchr/testify/assert"
)

const testLayout = "2006-01-02"

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
	err := storage.Remove(domain.EventID(uuid.New()))
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_RemoveIncorrectId(t *testing.T) {
	storage := New()
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	err := storage.Remove(domain.EventID(uuid.New()))
	assert.Error(t, err)
	assert.Equal(t, ErrNotExist, err)
}

func TestMemoryStorage_Remove(t *testing.T) {
	storage := New()

	event := domain.Event{ID: domain.EventID(uuid.New())}

	storage.Insert(event)
	storage.Insert(domain.Event{ID: domain.EventID(uuid.New())})

	err := storage.Remove(event.ID)
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

func TestMemoryStorage_GetEventsInTime(t *testing.T) {
	storage := New()

	getTime := func(date string) time.Time {
		res, _ := time.Parse(testLayout, date)
		return res
	}

	events := []domain.Event{
		{
			ID:       domain.EventID(uuid.New()),
			DateTime: getTime("2020-01-05"),
		},
		{
			ID:       domain.EventID(uuid.New()),
			DateTime: getTime("2020-01-10"),
		},
		{
			ID:       domain.EventID(uuid.New()),
			DateTime: getTime("2020-02-08"),
		},
	}

	for _, event := range events {
		storage.Insert(event)
	}

	result, err := storage.GetEventsInTime(getTime("2020-03-10"), time.Duration(time.Hour*24*10))

	assert.NoError(t, err)
	assert.Empty(t, result)

	result, err = storage.GetEventsInTime(getTime("2020-01-05"), time.Duration(time.Hour*24*40))

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	result, err = storage.GetEventsInTime(getTime("2020-01-06"), time.Duration(time.Hour*24*40))

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = storage.GetEventsInTime(getTime("2020-01-09"), time.Duration(time.Hour*24*1))

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	result, err = storage.GetEventsInTime(getTime("2020-01-07"), time.Duration(time.Hour*24*1))

	assert.NoError(t, err)
	assert.Empty(t, result)
}
