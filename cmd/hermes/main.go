package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api"
	"github.com/rugwirobaker/hermes/sqlite"
	"github.com/rugwirobaker/hermes/tracing"
	"go.opentelemetry.io/otel"
)

func main() {

	port := os.Getenv("PORT")
	id := os.Getenv("HELMES_SMS_APP_ID")
	secret := os.Getenv("HELMES_SMS_APP_SECRET")
	sender := os.Getenv("HELMES_SENDER_IDENTITY")
	callback := os.Getenv("HELMES_CALLBACK_URL")
	// uptraceDSN := os.Getenv("UPTRACE_DSN")
	dbURL := os.Getenv("DATABASE_URL")
	honeyCombKey := os.Getenv("HONEYCOMB_API_KEY")
	serviceName := os.Getenv("SERVICE_NAME")
	honeyCombDns := os.Getenv("HONEYCOMB_DSN")
	dbDriver := os.Getenv("DATABASE_DRIVER")

	if dbURL == "" {
		dbURL = "hermes.db"
	}

	if dbDriver == "" {
		dbDriver = "sqlite3"
	}

	ctx := context.Background()

	provider, err := tracing.Provider(ctx, honeyCombKey, honeyCombDns, serviceName)

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		_ = provider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(provider)

	cli := provideClient(provider)

	if strings.Contains(dbDriver, "sqlite") {
		dbURL = fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", dbURL)
	}

	db, err := sqlite.NewDB(dbURL, dbDriver, serviceName, provider)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := hermes.NewStore(db)

	service, err := hermes.NewSendService(cli, id, secret, sender, callback)
	if err != nil {
		log.Fatalf("could not initialize sms service: %v", err)
	}

	events := hermes.NewPubsub()
	defer events.Close()

	keys := hermes.NewIdempotencyKeyStore(db)

	log.Println("initialized hermes api")
	api := api.New(service, events, store, keys, provider)
	mux := chi.NewMux()
	mux.Mount("/api", api.Handler())

	if len(port) == 0 {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("starting application at port %v", port)

	err = srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-idleConnsClosed
}
