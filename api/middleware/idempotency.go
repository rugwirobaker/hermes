package middleware

import (
	"net/http"

	"github.com/rugwirobaker/hermes/api/request"
	"github.com/rugwirobaker/hermes/rand"
)

func Idempotency(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get(IdempotencyKeyHeader)

		if key == "" {
			key = rand.String(32, nil)
		}

		ctx := request.WithIdempotencyKey(r.Context(), key)
		r = r.WithContext(ctx)

		w.Header().Set(IdempotencyKeyHeader, key)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
