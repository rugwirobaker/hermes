package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rugwirobaker/hermes/build"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// start a new span for the request and adds a few iniial attributes
// request.method
// request.path
// request.host
// request.remote_addr
// request.user_agent
// request.idempotency_key
// request.id
// request.content_length
// request.content_type
func Tracing(provider trace.TracerProvider) Middleware {
	var tracer trace.Tracer

	if provider != nil {
		tracer = provider.Tracer(build.Info().ServiceName)
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w = &StatusWriter{ResponseWriter: w}

			ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
			defer span.End()
			r = r.WithContext(ctx)

			sw, ok := w.(*StatusWriter)
			if !ok {
				panic(fmt.Sprintf("ResponseWriter not a *tracing.StatusWriter; got %T", w))
			}

			// pass the span through the request context and serve the request to the next middleware
			next.ServeHTTP(sw, r)
			// capture response data
			EndHTTPSpan(r, sw.Status, span)

		}
		return http.HandlerFunc(fn)
	}
}

// EndHTTPSpan captures request and response data after the handler is done.
func EndHTTPSpan(r *http.Request, status int, span trace.Span) {
	// set the resource name as we get it only once the handler is executed
	route := chi.RouteContext(r.Context()).RoutePattern()
	span.SetName(fmt.Sprintf("%s %s", r.Method, route))
	span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...)
	span.SetAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...)
	span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", route, r)...)
	span.SetAttributes(semconv.HTTPRouteKey.String(route))

	// 0 status means one has not yet been sent in which case net/http library will write StatusOK
	if status == 0 {
		status = http.StatusOK
	}
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))
	span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(status, trace.SpanKindServer))
}
