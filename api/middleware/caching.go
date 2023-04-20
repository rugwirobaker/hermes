package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/request"
)

const (
	IdempotencyKeyHeader = "Idempotency-Key"
)

var cacheAbleMethods = []string{
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

func isCacheableMethod(method string, cacheAbleMethods []string) bool {
	for _, m := range cacheAbleMethods {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}

// Caching records(httptest.ResponseRecorder) and whole responses(http.Response) to the cache(hermes.Cache) with the Idempotency-Key header as the key.
// Next time it sees the same Idempotency-Key in a request, it will return the recorded response. If it sees a different Idempotency-Key, it will call the next handler.
func Caching(cache hermes.IdempotencyKeyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !isCacheableMethod(r.Method, cacheAbleMethods) {
				next.ServeHTTP(w, r)
				return
			}

			key, ok := request.IdempotencyFrom(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			entry, err := cache.Get(r.Context(), key)
			if err == nil && entry.Path == r.URL.Path {
				log.Printf("[Caching] Cache hit for key: %s", key)

				w.WriteHeader(entry.Code)
				for k, v := range entry.Headers {
					w.Header()[k] = v
				}

				w.Write(entry.Body)
				return
			}

			if err != hermes.ErrNotFound {
				log.Printf("[Caching] Cache error for key: %s, error: %v", key, err)
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("[Caching] Cache miss for key: %s", key)

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			record := &hermes.IdempotencyRecord{
				Key:     key,
				Code:    rec.Code,
				Headers: rec.Header().Clone(),
				Body:    rec.Body.Bytes(),
				Path:    r.URL.Path,
			}
			if err := cache.Set(r.Context(), record); err != nil {
				log.Printf("[Caching] Cache set error for key: %s, error: %v", key, err)
			}

			for k, v := range rec.Header() {
				w.Header()[k] = v
			}

			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
		}

		return http.HandlerFunc(fn)
	}
}
