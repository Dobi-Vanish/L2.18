package service

import (
	"l2.18/internal/model"
	"l2.18/internal/repository"
	"l2.18/pkg/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventService_CreateEvent_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Test Event",
	}

	err := service.CreateEvent(event)
	require.NoError(t, err)
	assert.NotZero(t, event.ID)
}

func TestEventService_CreateEvent_ValidationErrors(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	tests := []struct {
		name  string
		event *model.Event
		want  string
	}{
		{
			name:  "empty text",
			event: &model.Event{UserID: 1, Date: time.Now(), Text: ""},
			want:  "event text cannot be empty",
		},
		{
			name:  "zero user id",
			event: &model.Event{UserID: 0, Date: time.Now(), Text: "Test"},
			want:  "user ID is required",
		},
		{
			name:  "zero date",
			event: &model.Event{UserID: 1, Date: time.Time{}, Text: "Test"},
			want:  "event date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateEvent(tt.event)
			require.Error(t, err)

			validationErr, ok := err.(errors.ValidationError)
			require.True(t, ok, "Expected ValidationError, got %T", err)
			assert.Contains(t, validationErr.Error(), tt.want)
		})
	}
}

func TestEventService_UpdateEvent_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Original Text",
	}
	err := service.CreateEvent(event)
	require.NoError(t, err)

	updatedEvent := &model.Event{
		ID:     event.ID,
		UserID: 1,
		Date:   time.Now(),
		Text:   "Updated Text",
	}
	err = service.UpdateEvent(updatedEvent)
	require.NoError(t, err)
}

func TestEventService_UpdateEvent_ValidationErrors(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	tests := []struct {
		name  string
		event *model.Event
		want  string
	}{
		{
			name:  "zero id",
			event: &model.Event{ID: 0, UserID: 1, Date: time.Now(), Text: "Test"},
			want:  "event ID is required",
		},
		{
			name:  "empty text",
			event: &model.Event{ID: 1, UserID: 1, Date: time.Now(), Text: ""},
			want:  "event text cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateEvent(tt.event)
			require.Error(t, err)

			validationErr, ok := err.(errors.ValidationError)
			require.True(t, ok, "Expected ValidationError, got %T", err)
			assert.Contains(t, validationErr.Error(), tt.want)
		})
	}
}

func TestEventService_UpdateEvent_NotFound(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	event := &model.Event{
		ID:     999,
		UserID: 1,
		Date:   time.Now(),
		Text:   "Test Event",
	}

	err := service.UpdateEvent(event)
	require.Error(t, err)

	businessErr, ok := err.(errors.BusinessError)
	require.True(t, ok, "Expected BusinessError, got %T", err)
	assert.Contains(t, businessErr.Error(), "event not found")
}

func TestEventService_DeleteEvent_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	event := &model.Event{
		UserID: 1,
		Date:   time.Now(),
		Text:   "Test Event",
	}
	err := service.CreateEvent(event)
	require.NoError(t, err)

	err = service.DeleteEvent(event.ID)
	require.NoError(t, err)
}

func TestEventService_DeleteEvent_ValidationError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	err := service.DeleteEvent(0)
	require.Error(t, err)

	validationErr, ok := err.(errors.ValidationError)
	require.True(t, ok, "Expected ValidationError, got %T", err)
	assert.Contains(t, validationErr.Error(), "event ID is required")
}

func TestEventService_DeleteEvent_NotFound(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	err := service.DeleteEvent(999)
	require.Error(t, err)

	businessErr, ok := err.(errors.BusinessError)
	require.True(t, ok, "Expected BusinessError, got %T", err)
	assert.Contains(t, businessErr.Error(), "event not found")
}

func TestEventService_GetEventsDay(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	event1 := &model.Event{UserID: 1, Date: date, Text: "Event 1"}
	event2 := &model.Event{UserID: 1, Date: date.Add(2 * time.Hour), Text: "Event 2"}

	service.CreateEvent(event1)
	service.CreateEvent(event2)

	events, err := service.GetEventsDay(date)
	require.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestEventService_GetEventsWeek(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	monday := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday
	tuesday := monday.Add(24 * time.Hour)
	nextMonday := monday.Add(7 * 24 * time.Hour)

	event1 := &model.Event{UserID: 1, Date: monday, Text: "Monday Event"}
	event2 := &model.Event{UserID: 1, Date: tuesday, Text: "Tuesday Event"}

	service.CreateEvent(event1)
	service.CreateEvent(event2)

	events, err := service.GetEventsWeek(monday, nextMonday)
	require.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestEventService_GetEventsMonth(t *testing.T) {
	repo := repository.NewMemoryRepository()
	service := NewEventService(repo)

	monthStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	event1 := &model.Event{UserID: 1, Date: monthStart, Text: "Month Start Event"}
	event2 := &model.Event{UserID: 1, Date: monthEnd, Text: "Month End Event"}

	service.CreateEvent(event1)
	service.CreateEvent(event2)

	events, err := service.GetEventsMonth(monthStart, monthEnd)
	require.NoError(t, err)
	assert.Len(t, events, 2)
}
