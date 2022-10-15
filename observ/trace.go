package observ

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, op string) (context.Context, trace.Span) {
	return otel.Tracer("hermes").Start(ctx, op)
}
