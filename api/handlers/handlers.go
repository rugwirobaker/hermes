package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	helmes "github.com/rugwirobaker/helmes"
)

var startTime = time.Now()

// SendHandler ...
func SendHandler(svc helmes.SendService) http.HandlerFunc {
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

// VersionHandler handles version reporting
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := helmes.Data()
		JSON(w, res, http.StatusOK)
	}
}

// HealthHandler handles application health reporting
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

// JSON responds with json
func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
