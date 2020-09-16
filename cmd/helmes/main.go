package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rugwirobaker/helmes"
	"github.com/rugwirobaker/helmes/api"
)

func main() {
	var envfile string

	ctx := context.Background()

	flag.StringVar(&envfile, "env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	id := os.Getenv("HELMES_SMS_APP_ID")
	secret := os.Getenv("HELMES_SMS_APP_SECRET")
	sender := os.Getenv("HELMES_SENDER_IDENTITY")
	callback := os.Getenv("HELMES_CALLBACK_URL")
	// log.Printf("env: port-->%s", port)
	// log.Printf("env: helmes id-->%s", id)
	// log.Printf("env: helmes secret-->%s", secret)
	// log.Printf("env: helmes sender identity-->%s", sender)
	// log.Printf("env: helmes callback url-->%s", callback)

	cli := provideClient()
	if err != nil {
		log.Fatalf("could not initialize client: %v", err)
	}

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
