package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api"
	"github.com/rugwirobaker/hermes/api/middleware"
	"github.com/rugwirobaker/hermes/sqlite"
	"github.com/rugwirobaker/hermes/tracing"
	"go.opentelemetry.io/otel"
)

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

func runServe(ctx context.Context, args []string) (err error) {
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := newConfig()

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
		log.Printf("starting application at port %v", config.port)

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	<-signalCh
	log.Println("received signal, shutting down")

	// wait for server to shutdown with timeout
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		return fmt.Errorf("could not shutdown server: %v", err)
	}
	close(idleConnsClosed)

	<-idleConnsClosed

	log.Println("server shutdown")
	return

}
