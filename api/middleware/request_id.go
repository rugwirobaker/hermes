package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rugwirobaker/hermes/api/request"
)

// WithRequestID ensures a request id is in the
// request context by either the incoming header
// or creating a new one
func WithRequestID(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {
			requestID = r.Header.Get("Fly-Request-Id")
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := request.WithRequestID(r.Context(), requestID)

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
