package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

// Rate limiter. Allows 1 request per second with a burst of 5.
var limiter = rate.NewLimiter(1, 5)

// Healthcheck server.
type Server struct {
	Bot *Bot
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

// Create a new healthcheck server.
func NewServer(bot *Bot) *Server {
	return &Server{Bot: bot}
}

// Create a new healthcheck response.
func NewHealthCheckResponse(status string) *HealthCheckResponse {
	return &HealthCheckResponse{Status: status}
}

// Start Discord bot healthcheck server.
func (server *Server) StartServer() {
	http.HandleFunc("/status", server.rateLimit(server.healthCheckHandler))

	port := os.Getenv("PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// Healthcheck handler.
func (server *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := NewHealthCheckResponse(server.Bot.Status)

	w.Header().Set("Content-Type", "application/json")
	if server.Bot.Status == "ACTIVE" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// Rate limiter middleware.
func (server *Server) rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
