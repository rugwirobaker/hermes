package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	hermes "github.com/rugwirobaker/hermes"
)

var startTime = time.Now()

// SendHandler ...
func SendHandler(svc hermes.SendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		in := new(hermes.SMS)

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
