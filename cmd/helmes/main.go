package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/rugwirobaker/helmes"
	"github.com/rugwirobaker/helmes/api"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	id := os.Getenv("HELMES_SMS_APP_ID")
	secret := os.Getenv("HELMES_SMS_APP_SECRET")
	sender := os.Getenv("HELMES_SENDER_IDENTITY")
	callback := os.Getenv("HELMES_CALLBACK_URL")

	cli := provideClient()

	service, err := helmes.New(cli, id, secret, sender, callback)
	if err != nil {
		log.Fatalf("could not initialize sms service: %v", err)
	}

	log.Println("initialized sms client...")

	api := api.New(service)
	mux := chi.NewMux()
	mux.Mount("/api", api.Handler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Started helmes Server at port %s", port)
	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
