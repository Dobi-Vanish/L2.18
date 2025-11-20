package repository

import (
	"errors"
	"l2.18/internal/model"
	"sort"
	"sync"
	"time"
)

// MemoryRepository struct holds events.
type MemoryRepository struct {
	mu     sync.RWMutex
	events map[int]*model.Event
	nextID int
}

// NewMemoryRepository creates new MemoryRepository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		events: make(map[int]*model.Event),
		nextID: 1,
	}
}

// New creates new Memory Repository.
func New() Repository {
	return NewMemoryRepository()
}

// CreateEvent adds new event to the map.
func (r *MemoryRepository) CreateEvent(event *model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event.ID = r.nextID
	r.events[event.ID] = &model.Event{
		ID:     event.ID,
		UserID: event.UserID,
		Date:   event.Date,
		Text:   event.Text,
	}
	r.nextID++

	return nil
}

// UpdateEvent updates event in the map
func (r *MemoryRepository) UpdateEvent(id int, event *model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return errors.New("event not found")
	}

	r.events[id] = &model.Event{
		ID:     id,
		UserID: event.UserID,
		Date:   event.Date,
		Text:   event.Text,
	}

	return nil
}

// DeleteEvent deletes event from a map.
func (r *MemoryRepository) DeleteEvent(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return errors.New("event not found")
	}

	delete(r.events, id)
	return nil
}

// GetEventDay gets events for a provided day.
func (r *MemoryRepository) GetEventDay(date time.Time) ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*model.Event
	for _, event := range r.events {
		if isSameDay(event.Date, date) {
			events = append(events, event)
		}
	}

	sortEventsByDate(events)
	return events, nil
}

// GetEventWeek gets events for a week.
func (r *MemoryRepository) GetEventWeek(startDate, endDate time.Time) ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*model.Event
	for _, event := range r.events {
		if (event.Date.Equal(startDate) || event.Date.After(startDate)) &&
			(event.Date.Equal(endDate) || event.Date.Before(endDate)) {
			events = append(events, event)
		}
	}

	sortEventsByDate(events)
	return events, nil
}

// GetEventMonth gets events for a month.
func (r *MemoryRepository) GetEventMonth(startDate, endDate time.Time) ([]*model.Event, error) {
	return r.GetEventWeek(startDate, endDate)
}

// isSameDay checks if provided days for a week/month is not the same dates
func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// sortEventsByDate sorts events by date.
func sortEventsByDate(events []*model.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})
}
