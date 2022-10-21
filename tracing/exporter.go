package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

//Provider initiate and retrive TracerProvider intance
func Provider(ctx context.Context, honeyCombKey, url, service string) (*trace.TracerProvider, error) {

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
	)

	// Create a new OTLP exporter
	exporter, err := exporter(ctx, honeyCombKey, url)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tracerProvider)
	//For ignoring otel error
	// otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {}))
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider, nil
}

// exporter initiate exporter
func exporter(ctx context.Context, key, dns string) (*otlptrace.Exporter, error) {

	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(dns),
		otlptracehttp.WithHeaders(map[string]string{
			"x-honeycomb-team": key,
			"User-Agent":       "hermes/lumo",
		}),
	))

	if err != nil {
		return nil, fmt.Errorf("create otlp exporter: %w", err)
	}

	return exporter, nil
}
