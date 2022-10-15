package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/observ"
)

var startTime = time.Now()

// SendHandler ...
func SendHandler(svc hermes.SendService, store hermes.Store) http.HandlerFunc {
	const op = "handlers.SendHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		in := new(hermes.SMS)

		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Printf("failed to send sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}
		out, err := svc.Send(r.Context(), in)
		if err != nil {
			log.Printf("failed to send sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		msg := &hermes.Message{
			ProviderID: out.ID,
			Recipient:  in.Recipient,
			Payload:    in.Payload,
			Status:     "pending",
			Cost:       out.Cost,
		}

		if _, err := store.Insert(ctx, msg); err != nil {
			log.Printf("failed to save sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("sent sms to '%s'", in.Recipient)
		JSON(w, out, 200)
	}
}

func GetMessage(store hermes.Store) http.HandlerFunc {
	const op = "handlers.GetMessage"

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		id := chi.URLParam(r, "id")

		msg, err := store.MessageByID(ctx, id)
		if err != nil {
			log.Printf("failed to get sms %v", err)
			span.RecordError(err)
			HttpError(w, err, 500)
			return
		}
		JSON(w, msg, 200)
	}
}
