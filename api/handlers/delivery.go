package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rugwirobaker/hermes"
)

// DeliveryHandler handles delivery callback reception
func DeliveryHandler(events hermes.Pubsub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		in := new(event)

		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Printf("failed to serialize request")
			JSON(w, NewError(err.Error()), 500)
			return
		}
		events.Publish(r.Context(), convertEvent(in))

		JSON(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}
