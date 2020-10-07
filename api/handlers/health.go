package handlers

import (
	"net/http"
	"runtime"
	"time"

	helmes "github.com/rugwirobaker/helmes"
)

// HealthHandler reports the health of the application
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &helmes.Health{
			GitRev:     helmes.Data().Version,
			Uptime:     time.Since(startTime).Seconds(),
			Goroutines: runtime.NumGoroutine(),
		}
		JSON(w, res, http.StatusOK)
	}
}
