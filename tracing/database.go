package tracing

import (
	"context"
	"fmt"
	"strings"

	"github.com/nhatthm/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

// DBTraceDriver driver will register a new tracing sql driver and return the driver name.
func DBTraceDriver(tp trace.TracerProvider, driver, dns, service string) (string, error) {

	if service == "" {
		service = driver + "-db-service"
	}

	// Register the otelsql wrapper for the provided database driver.
	driverName, err := otelsql.RegisterWithSource(
		driver,
		dns,
		otelsql.WithDefaultAttributes(
			semconv.ServiceNameKey.String(service),
		),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.WithSystem(semconv.DBSystemKey.String(driver)),
		otelsql.WithTracerProvider(tp),
		otelsql.WithSpanNameFormatter(formatDBSpan),
	)

	if err != nil {
		return "", fmt.Errorf("error registering database tracing driver: %w", err)
	}

	return driverName, nil
}

func formatDBSpan(ctx context.Context, op string) string {
	const qPrefix = "-- name: "
	q := otelsql.QueryFromContext(ctx)
	if q == "" || !strings.HasPrefix(q, qPrefix) {
		return strings.ToUpper(op)
	}

	s := strings.SplitN(strings.TrimPrefix(q, qPrefix), " ", 2)[0]
	return fmt.Sprintf("%s %s", strings.ToUpper(op), s)
}
