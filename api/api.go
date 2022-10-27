package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/handlers"
	mw "github.com/rugwirobaker/hermes/api/middleware"
	"go.opentelemetry.io/otel/trace"
)

// Server ...
type Server struct {
	events   hermes.Pubsub
	service  hermes.SendService
	apps     hermes.AppStore
	messages hermes.Store
	cache    mw.Cache
	provider trace.TracerProvider
}

// New api Server instance
func New(
	svc hermes.SendService,
	events hermes.Pubsub,
	apps hermes.AppStore,
	messages hermes.Store,
	cache mw.Cache,
	provider trace.TracerProvider,
) *Server {
	return &Server{service: svc, events: events, apps: apps, messages: messages, cache: cache, provider: provider}
}

// Handler returns an http.Handler
func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(mw.Tracing(s.provider))
	r.Use(mw.WithRequestID)
	r.Use(mw.Idempotency)
	r.Use(mw.Caching(s.cache))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hermes the messenger of the gods at your service😉"))
	})

	r.Get("/version", handlers.VersionHandler())
	r.Get("/healthz", handlers.HealthHandler())

	r.Route("/messages", func(r chi.Router) {
		r.Use(mw.Authenticate(s.apps))
		r.Post("/send", handlers.SendHandler(s.service, s.messages, s.apps))
		r.Get("/serial/{id}", handlers.GetMessageBySerialID(s.messages))
		r.Get("/{id}", handlers.GetMessageByProviderID(s.messages))
		r.Get("/events/{id}/status", handlers.SubscribeHandler(s.events))
	})

	r.Route("/apps", func(r chi.Router) {
		r.Post("/", handlers.RegisterApp(s.apps))
		r.Get("/", handlers.ListApps(s.apps))
	})
	r.HandleFunc("/delivery", handlers.DeliveryHandler(s.events, s.messages))

	return r
}
