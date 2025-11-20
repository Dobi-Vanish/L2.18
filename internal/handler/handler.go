package handler

import (
	"encoding/json"
	"l2.18/internal/model"
	"l2.18/internal/service"
	"l2.18/pkg/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// EventHandler contains service's events.
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler creates new copy of EventHandler.
func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

// handleError determines error types and returns correct ones.
func (h *EventHandler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case errors.ValidationError:
		http.Error(w, e.Error(), http.StatusBadRequest)
	case errors.BusinessError:
		http.Error(w, e.Error(), http.StatusServiceUnavailable)
	case errors.InternalError:
		http.Error(w, e.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CreateEvent handler to create event with provided info.
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		})
		return
	}

	if err := h.service.CreateEvent(&event); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// UpdateEvent updates event by id with provided info.
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "id",
			Message: "invalid event ID format",
		})
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		})
		return
	}
	event.ID = id

	if err := h.service.UpdateEvent(&event); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DeleteEvent deletes event by id.
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "id",
			Message: "invalid event ID format",
		})
		return
	}

	if err := h.service.DeleteEvent(id); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// GetEventsForDay gets all events for a day
func (h *EventHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	userIDStr := r.URL.Query().Get("user_id")

	if dateStr == "" || userIDStr == "" {
		h.handleError(w, errors.ValidationError{
			Field:   "query_params",
			Message: "date and user_id parameters are required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "date",
			Message: "invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "user_id",
			Message: "invalid user ID format",
		})
		return
	}

	events, err := h.service.GetEventsDay(date)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var userEvents []*model.Event
	for _, event := range events {
		if event.UserID == userID {
			userEvents = append(userEvents, event)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userEvents)
}

// GetEventsForWeek gets events for a week.
func (h *EventHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	userIDStr := r.URL.Query().Get("user_id")

	if dateStr == "" || userIDStr == "" {
		h.handleError(w, errors.ValidationError{
			Field:   "query_params",
			Message: "date and user_id parameters are required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "date",
			Message: "invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "user_id",
			Message: "invalid user ID format",
		})
		return
	}

	weekStart := date
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.AddDate(0, 0, -1)
	}
	weekEnd := weekStart.AddDate(0, 0, 7)

	events, err := h.service.GetEventsWeek(weekStart, weekEnd)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var userEvents []*model.Event
	for _, event := range events {
		if event.UserID == userID {
			userEvents = append(userEvents, event)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userEvents)
}

// GetEventsForMonth get events for a month.
func (h *EventHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	userIDStr := r.URL.Query().Get("user_id")

	if dateStr == "" || userIDStr == "" {
		h.handleError(w, errors.ValidationError{
			Field:   "query_params",
			Message: "date and user_id parameters are required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "date",
			Message: "invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.handleError(w, errors.ValidationError{
			Field:   "user_id",
			Message: "invalid user ID format",
		})
		return
	}

	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	events, err := h.service.GetEventsMonth(monthStart, monthEnd)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var userEvents []*model.Event
	for _, event := range events {
		if event.UserID == userID {
			userEvents = append(userEvents, event)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userEvents)
}
