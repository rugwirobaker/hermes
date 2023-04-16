module github.com/rugwirobaker/hermes

go 1.16

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.15
	github.com/nhatthm/otelsql v0.4.0
	github.com/rubenv/sql-migrate v1.2.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.3
	go.opentelemetry.io/otel v1.11.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.11.0
)
