package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/rugwirobaker/hermes/api/request"
)

const (
	IdempotencyKeyHeader = "Idempotency-Key"
)

// var cacheAbleMethods = []string{
// 	http.MethodPost,
// 	http.MethodPut,
// 	http.MethodDelete,
// }

// Caching records(httptest.ResponseRecorder) and whole responses(http.Response) to the cache(hermes.Cache) with the Idempotency-Key header as the key.
// Next time it sees the same Idempotency-Key in a request, it will return the recorded response. If it sees a different Idempotency-Key, it will call the next handler.
func Caching(cache Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			key, ok := request.IdempotencyFrom(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			if entry, ok := cache.Get(key); ok && entry.Path == r.URL.Path {
				log.Printf("[Caching] Cache hit for key: %s", key)

				w.WriteHeader(entry.Code)
				for k, v := range entry.Headers {
					w.Header()[k] = v
				}

				w.Write(entry.Body)
				return
			}

			log.Printf("[Caching] Cache miss for key: %s", key)

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)
			entry := NewEntry(rec.Code, rec.Header(), r.URL.Path, rec.Body.Bytes())
			cache.Set(key, entry)
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
		}

		return http.HandlerFunc(fn)
	}
}

type Cache interface {
	Get(key string) (*Entry, bool)
	Set(key string, entry *Entry)
}

type Entry struct {
	Code    int
	Headers http.Header
	Body    []byte
	Path    string
}

func NewEntry(code int, headers http.Header, path string, body []byte) *Entry {
	return &Entry{
		Code:    code,
		Headers: headers,
		Body:    body,
		Path:    path,
	}
}

type memoryCache struct {
	mu    *sync.RWMutex
	cache map[string]*Entry
}

func NewMemoryCache() Cache {
	return &memoryCache{
		mu:    &sync.RWMutex{},
		cache: make(map[string]*Entry),
	}
}

func (c *memoryCache) Get(key string) (*Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	return entry, ok
}

func (c *memoryCache) Set(key string, entry *Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = entry
}

// func isCachable(r *http.Request) bool {
// 	for _, method := range cacheAbleMethods {
// 		if method == r.Method {
// 			return true
// 		}
// 	}
// 	return false
// }
