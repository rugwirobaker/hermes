package handlers

import (
	"log"
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/render"
	"github.com/rugwirobaker/hermes/api/request"
	"github.com/rugwirobaker/hermes/observ"
)

// DeliveryHandler handles delivery callback reception
func DeliveryHandler(events hermes.Pubsub, store hermes.Store) http.HandlerFunc {
	const op = "handlers.DeliveryHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		_, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		if r.Method != http.MethodPost {
			err := hermes.NewErrInvalid("invalid method")
			span.RecordError(err)
			render.HttpError(w, err)
			return
		}

		var in hermes.Callback

		if err := request.Decode(r.Context(), r.Body, &in); err != nil {
			log.Printf("failed to decode request body %v", err)
			span.RecordError(err)
			render.HttpError(w, err)
			return
		}

		// find message in store
		msg, err := store.MessageByID(r.Context(), in.MsgRef)
		if err != nil {
			log.Printf("failed to find message in store")
			span.RecordError(err)
			render.HttpError(w, err)
			return
		}

		// update message status
		msg.Status = hermes.St(in.Status)

		if _, err := store.Update(r.Context(), msg); err != nil {
			log.Printf("failed to update message status")
			span.RecordError(err)
			render.HttpError(w, err)
			return
		}

		events.Publish(r.Context(), convertEvent(in))

		render.JSON(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}

func convertEvent(in hermes.Callback) hermes.Event {
	return hermes.Event{
		ID:        in.MsgRef,
		Status:    hermes.St(in.Status),
		Recipient: in.Recipient,
	}
}
