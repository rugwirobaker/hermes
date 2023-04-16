package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	hermes "github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/request"
	"github.com/rugwirobaker/hermes/observ"
)

var startTime = time.Now()

// SendHandler ...
func SendHandler(svc hermes.SendService, messages hermes.Store, apps hermes.AppStore) http.HandlerFunc {
	const op = "handlers.SendHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		app, ok := request.AppFrom(ctx)
		if !ok {
			span.RecordError(hermes.ErrUnauthorized)
			HttpError(w, hermes.ErrUnauthorized, 401)
			return
		}

		in := new(hermes.SMS)

		// call Decode to decode the request body into the struct
		if err := request.Decode(ctx, r.Body, in); err != nil {
			log.Printf("failed to send sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		in.Sender = app.Sender

		out, err := svc.Send(ctx, in)
		if err != nil {
			log.Printf("failed to send sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		msg := &hermes.Message{
			ProviderID: out.ID,
			From:       app.ID,
			Recipient:  in.Recipient,
			Payload:    in.Payload,
			Status:     hermes.Status(out.Status),
			Cost:       out.Cost,
		}

		if _, err := messages.Insert(ctx, msg); err != nil {
			log.Printf("failed to save sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		// update the sent message count for the app
		app.MessageCount = app.MessageCount + out.Count

		if err := apps.Update(ctx, app); err != nil {
			log.Printf("failed to update app %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		if span.IsRecording() {
			span.SetAttributes(
				observ.String("message.id", out.ID),
				observ.Int64("message.serial_id", int64(msg.ID)),
				observ.String("message.cost", fmt.Sprintf("%f", out.Cost)),
			)
		}

		log.Printf("sent sms to '%s'", in.Recipient)
		JSON(w, out, 200)
	}
}

func GetMessageBySerialID(store hermes.Store) http.HandlerFunc {
	const op = "handlers.GetMessage"

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		id := chi.URLParam(r, "id")

		msg, err := store.MessageBySerial(ctx, id)
		if err != nil {
			log.Printf("failed to get sms %v", err)
			span.RecordError(err)
			HttpError(w, err, 500)
			return
		}
		JSON(w, msg, 200)
	}
}

func GetMessageByProviderID(store hermes.Store) http.HandlerFunc {
	const op = "handlers.GetMessageBYProviderID"

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
