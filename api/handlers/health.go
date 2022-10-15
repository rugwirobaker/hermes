package handlers

import (
	"net/http"
	"os"
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
			Region:     os.Getenv("FLY_REGION"),
		}
		JSON(w, res, http.StatusOK)
	}
}
