package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api"
	"github.com/rugwirobaker/hermes/api/middleware"
	"github.com/rugwirobaker/hermes/build"
	"github.com/rugwirobaker/hermes/sqlite"
	"github.com/rugwirobaker/hermes/tracing"
	"go.opentelemetry.io/otel"
)

func main() {
	config := newConfig()

	log.Printf("database: %s", config.dsn)

	ctx := context.Background()

	provider, err := tracing.Provider(
		ctx,
		config.honecomb.apiKey,
		config.honecomb.dsn,
		config.serviceName,
		config.environment,
		config.region,
		config.hostID,
	)

	if err != nil {
		log.Fatalf("could not initialize tracing provider: %v", err)
	}

	defer func() {
		_ = provider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(provider)

	client := createSmsClient(config.apiKey, provider)

	db, err := sqlite.NewDB(config.dsn, config.serviceName, provider)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	apps := hermes.NewAppStore(db)

	messages := hermes.NewStore(db)

	service, err := hermes.NewSendService(client, config.sender)

	if err != nil {
		log.Fatalf("could not initialize sms service: %v", err)
	}

	events := hermes.NewPubsub()
	defer events.Close()

	cache := middleware.NewMemoryCache()

	log.Println("initialized hermes api")
	api := api.New(service, events, apps, messages, cache, provider)
	mux := chi.NewMux()

	// only attach mw.FlyReplay if we're running on fly.io
	if config.flyAppName != "" {
		mux.Use(middleware.FlyReplay(config.dsn))
	}

	mux.Mount("/api", api.Handler())

	srv := NewServer(config.port, mux)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("starting application at port %v", config.port)

	err = srv.Start()

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-idleConnsClosed
}

type config struct {
	port     string
	apiKey   string
	sender   string
	dsn      string
	honecomb struct {
		apiKey string
		dsn    string
	}
	flyAppName  string
	serviceName string
	environment string
	region      string
	hostID      string
}

func newConfig() *config {

	config := &config{
		port:        os.Getenv("PORT"),
		apiKey:      os.Getenv("PINDO_API_KEY"),
		sender:      os.Getenv("HELMES_SENDER_IDENTITY"),
		dsn:         os.Getenv("DATABASE_URL"),
		serviceName: build.Info().ServiceName,
		environment: "development",
		flyAppName:  os.Getenv("FLY_APP_NAME"),
		region:      "local",
		hostID:      "local",
		honecomb: struct {
			apiKey string
			dsn    string
		}{
			apiKey: os.Getenv("HONEYCOMB_API_KEY"),
			dsn:    os.Getenv("HONEYCOMB_DSN"),
		},
	}

	if os.Getenv("ENVIRONMENT") != "" {
		config.environment = os.Getenv("ENVIRONMENT")
	}

	if os.Getenv("FLY_REGION") != "" {
		config.region = os.Getenv("FLY_REGION")
	}

	if config.port == "" {
		config.port = "8080"
	}

	var err error

	if config.hostID, err = os.Hostname(); err != nil {
		log.Printf("warn:unable to get hostname: %v", err)
	}

	if os.Getenv("FLY_ALLOC_ID") != "" {
		config.hostID = os.Getenv("FLY_ALLOC_ID")
	}

	if config.dsn == "" {
		config.dsn = "hermes.db"
	}
	return config
}

// this is to make it easy to run several services in the same process on different ports
// metrics for example should be on a different port
type Server struct {
	*http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:        ":" + port,
			Handler:     handler,
			IdleTimeout: 5 * time.Second,
			ReadTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
