package main

import (
	"encoding/json"
	"l2.18/internal/config"
	"l2.18/internal/handler"
	"l2.18/internal/repository"
	"l2.18/internal/service"
	"l2.18/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()
	repo := repository.New()

	eventService := service.NewEventService(repo)
	eventHandler := handler.NewEventHandler(eventService)

	router := mux.NewRouter()
	eventHandler.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	handlerWithMiddleware := middleware.LoggingMiddleware("logs/requests.log")(router)

	log.Printf("Server starting on: 8081")
	log.Printf("File logging enabled: logs/requests.log")
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPServerPort, handlerWithMiddleware))
}
