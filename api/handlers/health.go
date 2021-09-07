package handlers

import (
	"net/http"
	"runtime"
	"time"

	hermes "github.com/rugwirobaker/hermes"
)

// HealthHandler reports the health of the application
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &hermes.Health{
			GitRev:     hermes.Data().Version,
			Uptime:     time.Since(startTime).Seconds(),
			Goroutines: runtime.NumGoroutine(),
		}
		JSON(w, res, http.StatusOK)
	}
}
