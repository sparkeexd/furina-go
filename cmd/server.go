package main

import (
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

// Rate limiter. Allows 1 request per second with a burst of 5.
var limiter = rate.NewLimiter(1, 5)

// Healthcheck server.
type Server struct {
	Bot *Bot
}

// Create a new healthcheck server.
func NewServer(bot *Bot) *Server {
	return &Server{Bot: bot}
}

// Start Discord bot healthcheck server.
func (server *Server) StartServer() {
	http.HandleFunc("/health-check", server.rateLimit(server.healthCheckHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// Health check handler.
func (server *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if server.Bot.Status == "Active" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "Active"}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "Inactive"}`))
	}
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
