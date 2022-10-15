package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/handlers"
	mw "github.com/rugwirobaker/hermes/api/middleware"
	"go.opentelemetry.io/otel/trace"
)

// Server ...
type Server struct {
	events   hermes.Pubsub
	service  hermes.SendService
	store    hermes.Store
	cache    mw.Cache
	provider trace.TracerProvider
}

// New api Server instance
func New(svc hermes.SendService, events hermes.Pubsub, store hermes.Store, cache mw.Cache, provider trace.TracerProvider) *Server {
	return &Server{service: svc, events: events, store: store, cache: cache, provider: provider}
}

// Handler returns an http.Handler
func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(otelchi.Middleware("hermes", otelchi.WithTracerProvider(s.provider)))
	r.Use(mw.Idempotency)
	r.Use(mw.Caching(s.cache))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to hermes"))
	})

	r.Get("/version", handlers.VersionHandler())
	r.Get("/healthz", handlers.HealthHandler())
	r.Post("/send", handlers.SendHandler(s.service, s.store))
	r.Get("/message/{id}", handlers.GetMessage(s.store))
	r.Get("/events/{id}/status", handlers.SubscribeHandler(s.events))
	r.HandleFunc("/delivery", handlers.DeliveryHandler(s.events))

	return r
}
