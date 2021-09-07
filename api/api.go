package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/handlers"
)

// Server ...
type Server struct {
	Events  hermes.Pubsub
	Service hermes.SendService
}

// New api Server instance
func New(svc hermes.SendService, events hermes.Pubsub) *Server {
	return &Server{Service: svc, Events: events}
}

// Handler returns an http.Handler
func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to hermes"))
	})

	r.Get("/version", handlers.VersionHandler())
	r.Get("/healthz", handlers.HealthHandler())
	r.Post("/send", handlers.SendHandler(s.Service))
	r.Get("/events/{id}/status", handlers.SubscribeHandler(s.Events))
	r.Post("/delivery", handlers.DeliveryHandler(s.Events))

	return r
}
