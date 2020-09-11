package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rugwirobaker/sam"
	"github.com/rugwirobaker/sam/api"
)

func main() {
	//port
	//sms-app-id
	//sms-app-secrert
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

	log.Printf("starting sam server at %s", port)

	http.ListenAndServe(":"+port, mux)
}
