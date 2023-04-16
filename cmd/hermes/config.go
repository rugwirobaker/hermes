package main

import (
	"log"
	"os"

	"github.com/rugwirobaker/hermes/build"
)

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
