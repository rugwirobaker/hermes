package handlers

import (
	"log"
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/request"
	"github.com/rugwirobaker/hermes/observ"
)

func RegisterApp(store hermes.AppStore) http.HandlerFunc {
	const op = "handlers.RegisterAppHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		in := new(hermes.App)

		// call Decode to decode the request body into the struct
		if err := request.Decode(ctx, r.Body, in); err != nil {
			log.Printf("failed to send sms %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		// do some validation
		if in.Name == "" {
			http.Error(w, "app name is required", 400)
			return
		}

		if in.Sender == "" {
			http.Error(w, "app sender is required", 400)
			return
		}

		// generate random api key
		apiKey, err := hermes.RandomString(18)
		if err != nil {
			log.Printf("failed to generate api key %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}
		in.APIKey = apiKey

		if err := store.Register(ctx, in); err != nil {
			log.Printf("failed to register app %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		if span.IsRecording() {
			span.SetAttributes(
				observ.String("app.name", in.Name),
			)
		}

		JSON(w, in, http.StatusOK)
	}
}

func ListApps(store hermes.AppStore) http.HandlerFunc {
	const op = "handlers.ListAppsHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		apps, err := store.List(ctx)
		if err != nil {
			log.Printf("failed to list apps %v", err)
			span.RecordError(err)
			http.Error(w, err.Error(), 500)
			return
		}

		if span.IsRecording() {
			span.SetAttributes(
				observ.Int64("retrieved.count", int64(len(apps))),
			)
		}

		JSON(w, apps, http.StatusOK)
	}
}
