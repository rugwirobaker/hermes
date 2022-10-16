package observ

import (
	"context"

	"github.com/rugwirobaker/hermes/build"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, op string) (context.Context, trace.Span) {
	return otel.Tracer(build.Info().ServiceName).Start(ctx, op)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func String(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

func Int64(key string, value int64) attribute.KeyValue {
	return attribute.Int64(key, value)
}
