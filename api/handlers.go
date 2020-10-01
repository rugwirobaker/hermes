package api

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	helmes "github.com/rugwirobaker/helmes"
)

var startTime = time.Now()

// SMSHandler ...
func SMSHandler(svc helmes.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		in := new(helmes.SMS)

		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Printf("failed to send sms %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		out, err := svc.Send(r.Context(), in)
		if err != nil {
			log.Printf("failed to send sms %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("sent sms to '%s'", in.Recipient)
		JSON(w, out, 200)
	}
}

// VersionHandler ...
func VersionHandler(svc helmes.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		build, err := svc.Version(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		JSON(w, build, http.StatusOK)
	}
}

// HealthHandler reports the health of the application
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			GitRev     string  `json:"git_rev"`
			Uptime     float64 `json:"uptime"`
			Goroutines int     `json:"goroutines"`
		}{
			GitRev:     helmes.Data().Version,
			Uptime:     time.Since(startTime).Seconds(),
			Goroutines: runtime.NumGoroutine(),
		}
		JSON(w, res, http.StatusOK)
	}
}

// JSON responds with json
func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
