package repository

import (
	"l2.18/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_CreateEvent(t *testing.T) {
	repo := NewMemoryRepository()

	event := &model.Event{
		UserID: 1,
		Date:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Text:   "Test Event",
	}

	err := repo.CreateEvent(event)
	require.NoError(t, err)
	assert.NotZero(t, event.ID)
	assert.Equal(t, 1, event.ID)

	events, err := repo.GetEventDay(event.Date)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, event.ID, events[0].ID)
	assert.Equal(t, event.Text, events[0].Text)
}

func TestMemoryRepository_UpdateEvent(t *testing.T) {
	repo := NewMemoryRepository()

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Original Text",
	}
	err := repo.CreateEvent(event)
	require.NoError(t, err)

	updatedEvent := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Updated Text",
	}
	err = repo.UpdateEvent(event.ID, updatedEvent)
	require.NoError(t, err)

	events, err := repo.GetEventDay(event.Date)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "Updated Text", events[0].Text)
}

func TestMemoryRepository_UpdateEvent_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Test Event",
	}

	err := repo.UpdateEvent(999, event)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestMemoryRepository_DeleteEvent(t *testing.T) {
	repo := NewMemoryRepository()

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Test Event",
	}
	err := repo.CreateEvent(event)
	require.NoError(t, err)

	err = repo.DeleteEvent(event.ID)
	require.NoError(t, err)

	events, err := repo.GetEventDay(event.Date)
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestMemoryRepository_DeleteEvent_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.DeleteEvent(999)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestMemoryRepository_GetEventDay(t *testing.T) {
	repo := NewMemoryRepository()

	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	event1 := &model.Event{
		UserID: 1,
		Date:   date,
		Text:   "Morning Event",
	}
	event2 := &model.Event{
		UserID: 1,
		Date:   date.Add(2 * time.Hour),
		Text:   "Afternoon Event",
	}
	event3 := &model.Event{
		UserID: 1,
		Date:   date.Add(24 * time.Hour),
		Text:   "Next Day Event",
	}

	repo.CreateEvent(event1)
	repo.CreateEvent(event2)
	repo.CreateEvent(event3)

	events, err := repo.GetEventDay(date)
	require.NoError(t, err)
	assert.Len(t, events, 2)

	assert.Equal(t, "Morning Event", events[0].Text)
	assert.Equal(t, "Afternoon Event", events[1].Text)
}

func TestMemoryRepository_GetEventWeek(t *testing.T) {
	repo := NewMemoryRepository()

	monday := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC) // Monday
	tuesday := monday.Add(24 * time.Hour)
	nextMonday := monday.Add(7 * 24 * time.Hour)

	event1 := &model.Event{UserID: 1, Date: monday, Text: "Monday Event"}
	event2 := &model.Event{UserID: 1, Date: tuesday, Text: "Tuesday Event"}
	event3 := &model.Event{UserID: 1, Date: nextMonday, Text: "Next Monday Event"}

	repo.CreateEvent(event1)
	repo.CreateEvent(event2)
	repo.CreateEvent(event3)

	events, err := repo.GetEventWeek(monday, nextMonday)
	require.NoError(t, err)
	assert.Len(t, events, 3)
}

func TestMemoryRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryRepository()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			event := &model.Event{
				UserID: 1,
				Date:   time.Now(),
				Text:   string(rune(id)),
			}
			_ = repo.CreateEvent(event)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	events, err := repo.GetEventDay(time.Now())
	require.NoError(t, err)
	assert.Len(t, events, 10)
}
