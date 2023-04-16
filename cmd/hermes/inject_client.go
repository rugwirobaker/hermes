package main

import (
	"net/http"

	"github.com/rugwirobaker/hermes/pindo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func createSmsClient(apiKey string, tracerProvider trace.TracerProvider) *pindo.Client {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithTracerProvider(tracerProvider),
			otelhttp.WithPublicEndpoint(),
		),
	}

	return pindo.NewWithClient(apiKey, httpClient)
}
