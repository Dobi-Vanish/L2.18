package service

import (
	"l2.18/internal/model"
	"l2.18/internal/repository"
	"l2.18/pkg/errors"
	"time"
)

// EventService struct holds repository for events.
type EventService struct {
	repo repository.Repository
}

// NewEventService creates new EventService.
func NewEventService(repo repository.Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

// CreateEvent creates event.
func (s *EventService) CreateEvent(event *model.Event) error {
	if event.Text == "" {
		return errors.ValidationError{
			Field:   "text",
			Message: "event text cannot be empty",
		}
	}
	if event.UserID == 0 {
		return errors.ValidationError{
			Field:   "user_id",
			Message: "user ID is required",
		}
	}
	if event.Date.IsZero() {
		return errors.ValidationError{
			Field:   "date",
			Message: "event date is required",
		}
	}

	err := s.repo.CreateEvent(event)
	if err != nil {
		return errors.InternalError{
			Operation: "create_event",
			Message:   err.Error(),
		}
	}
	return nil
}

// UpdateEvent updates event by and with provided info.
func (s *EventService) UpdateEvent(event *model.Event) error {
	if event.ID == 0 {
		return errors.ValidationError{
			Field:   "id",
			Message: "event ID is required",
		}
	}
	if event.UserID == 0 {
		return errors.ValidationError{
			Field:   "user_id",
			Message: "user ID is required",
		}
	}
	if event.Text == "" {
		return errors.ValidationError{
			Field:   "text",
			Message: "event text cannot be empty",
		}
	}

	err := s.repo.UpdateEvent(event.ID, event)
	if err != nil {
		if err.Error() == "event not found" {
			return errors.BusinessError{
				Operation: "update_event",
				Message:   "event not found",
			}
		}
		return errors.InternalError{
			Operation: "update_event",
			Message:   err.Error(),
		}
	}
	return nil
}

// DeleteEvent deletes event by provided id.
func (s *EventService) DeleteEvent(eventID int) error {
	if eventID == 0 {
		return errors.ValidationError{
			Field:   "id",
			Message: "event ID is required",
		}
	}

	err := s.repo.DeleteEvent(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			return errors.BusinessError{
				Operation: "delete_event",
				Message:   "event not found",
			}
		}
		return errors.InternalError{
			Operation: "delete_event",
			Message:   err.Error(),
		}
	}
	return nil
}

// GetEventsDay get all events for a day.
func (s *EventService) GetEventsDay(date time.Time) ([]*model.Event, error) {
	events, err := s.repo.GetEventDay(date)
	if err != nil {
		return nil, errors.InternalError{
			Operation: "get_events_day",
			Message:   err.Error(),
		}
	}
	return events, nil
}

// GetEventsWeek get all events for a week.
func (s *EventService) GetEventsWeek(dayStart, dayEnd time.Time) ([]*model.Event, error) {
	events, err := s.repo.GetEventWeek(dayStart, dayEnd)
	if err != nil {
		return nil, errors.InternalError{
			Operation: "get_events_week",
			Message:   err.Error(),
		}
	}
	return events, nil
}

// GetEventsMonth get all events for a month.
func (s *EventService) GetEventsMonth(dayStart, dayEnd time.Time) ([]*model.Event, error) {
	events, err := s.repo.GetEventMonth(dayStart, dayEnd)
	if err != nil {
		return nil, errors.InternalError{
			Operation: "get_events_month",
			Message:   err.Error(),
		}
	}
	return events, nil
}
