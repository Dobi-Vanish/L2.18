package repository

import (
	"l2.18/internal/model"
	"time"
)

// Repository interface that holds function for CRUD operations with events.
type Repository interface {
	CreateEvent(event *model.Event) error
	UpdateEvent(id int, updateEvent *model.Event) error
	DeleteEvent(eventID int) error
	GetEventDay(date time.Time) ([]*model.Event, error)
	GetEventWeek(dateStart, dateEnd time.Time) ([]*model.Event, error)
	GetEventMonth(dateStart, dateEnd time.Time) ([]*model.Event, error)
}
