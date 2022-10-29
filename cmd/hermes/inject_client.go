package main

import (
	"net/http"

	"github.com/quarksgroup/sms-client/sms"
	"github.com/quarksgroup/sms-client/sms/driver/fdi"
	"github.com/quarksgroup/sms-client/sms/transport/oauth2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func provideClient(tracerProvider trace.TracerProvider) *sms.Client {
	client := fdi.NewDefault()
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Scheme: oauth2.SchemeBearer,
			Source: oauth2.ContextTokenSource(),
			Base: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithTracerProvider(tracerProvider),
				otelhttp.WithPublicEndpoint(),
			),
		},
	}
	return client
}
