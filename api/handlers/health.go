package handlers

import (
	"net/http"
	"os"
	"runtime"
	"time"

	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/observ"
)

// HealthHandler reports the health of the application
func HealthHandler() http.HandlerFunc {
	const op = "handlers.HealthHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		_, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		res := &hermes.Health{
			GitRev:     hermes.Data().Version,
			Uptime:     time.Since(startTime).Seconds(),
			Goroutines: runtime.NumGoroutine(),
			Region:     os.Getenv("FLY_REGION"),
		}
		JSON(w, res, http.StatusOK)
	}
}
