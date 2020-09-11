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
	"github.com/rugwirobaker/sam"
	"github.com/rugwirobaker/sam/api"
)

func main() {
	var envfile string
	flag.StringVar(&envfile, "env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	// id := os.Getenv("SAM_SMS_APP_ID")
	// secret := os.Getenv("SAM_SMS_APP_SECRET")

	service := sam.New()
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
	log.Printf("Started Sam Server at port %s", port)
	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
