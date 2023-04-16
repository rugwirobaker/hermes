package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rugwirobaker/hermes/sqlite"
	"github.com/rugwirobaker/hermes/tracing"
	"go.opentelemetry.io/otel"
)

func runMigrate(ctx context.Context, args []string) (err error) {
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

	db, err := sqlite.NewDB(config.dsn, config.serviceName, provider)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	n, err := db.Migrate(sqlite.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	fmt.Printf("applied %d migrations\n", n)

	return
}
