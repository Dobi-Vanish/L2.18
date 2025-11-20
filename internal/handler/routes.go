package handler

import "github.com/gorilla/mux"

// RegisterRoutes registers all available routes.
func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_event", h.CreateEvent).Methods("POST")
	router.HandleFunc("/update_event/{id}", h.UpdateEvent).Methods("POST")
	router.HandleFunc("/delete_event/{id}", h.DeleteEvent).Methods("POST")
	router.HandleFunc("/events_for_day", h.GetEventsForDay).Methods("GET")
	router.HandleFunc("/events_for_week", h.GetEventsForWeek).Methods("GET")
	router.HandleFunc("/events_for_month", h.GetEventsForMonth).Methods("GET")
}
