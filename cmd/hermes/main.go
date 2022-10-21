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
	"github.com/rugwirobaker/hermes/build"
	"github.com/rugwirobaker/hermes/sqlite"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
)

func main() {

	port := os.Getenv("PORT")
	id := os.Getenv("HELMES_SMS_APP_ID")
	secret := os.Getenv("HELMES_SMS_APP_SECRET")
	sender := os.Getenv("HELMES_SENDER_IDENTITY")
	callback := os.Getenv("HELMES_CALLBACK_URL")
	uptraceDSN := os.Getenv("UPTRACE_DSN")
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		dbURL = "hermes.db"
	}

	ctx := context.Background()

	info := build.Info()

	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		uptrace.WithDSN(uptraceDSN),
		uptrace.WithServiceName(info.ServiceName),
		uptrace.WithServiceVersion(info.Version),
	)

	provider := uptrace.TracerProvider()
	defer provider.Shutdown(ctx)

	otel.SetTracerProvider(provider)

	cli := provideClient(provider)

	db, err := sqlite.NewDB(dbURL, provider)
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
