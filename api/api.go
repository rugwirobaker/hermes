package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/handlers"
	mw "github.com/rugwirobaker/hermes/api/middleware"
)

// Server ...
type Server struct {
	Events  hermes.Pubsub
	Service hermes.SendService
	Cache   mw.Cache
}

// New api Server instance
func New(svc hermes.SendService, events hermes.Pubsub, cache mw.Cache) *Server {
	return &Server{Service: svc, Events: events, Cache: cache}
}

// Handler returns an http.Handler
func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.Idempotency)
	r.Use(mw.Caching(s.Cache))

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
