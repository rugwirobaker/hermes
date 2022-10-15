package handlers

import (
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/observ"
)

// VersionHandler ...
func VersionHandler() http.HandlerFunc {
	const op = "handlers.VersionHandler"

	return func(w http.ResponseWriter, r *http.Request) {
		_, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		JSON(w, hermes.Data(), http.StatusOK)
	}
}
