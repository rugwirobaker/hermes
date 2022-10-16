package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/observ"
)

// DeliveryHandler handles delivery callback reception
func DeliveryHandler(events hermes.Pubsub, store hermes.Store) http.HandlerFunc {
	const op = "handlers.DeliveryHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		_, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusOK)
			return
		}

		var in hermes.Callback

		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Printf("failed to serialize request")
			JSON(w, NewError(err.Error()), 500)
			return
		}

		// find message in store
		msg, err := store.MessageByID(r.Context(), in.MsgRef)
		if err != nil {
			log.Printf("failed to find message in store")
			JSON(w, NewError(err.Error()), 500)
			return
		}

		// update message status
		msg.Status = hermes.St(in.Status)

		if _, err := store.Update(r.Context(), msg); err != nil {
			log.Printf("failed to update message status")
			JSON(w, NewError(err.Error()), 500)
			return
		}

		events.Publish(r.Context(), convertEvent(in))

		JSON(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}

func convertEvent(in hermes.Callback) hermes.Event {
	return hermes.Event{
		ID:        in.MsgRef,
		Status:    hermes.St(in.Status),
		Recipient: in.Recipient,
	}
}
