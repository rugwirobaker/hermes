package handlers

import (
	"net/http"

	"github.com/rugwirobaker/hermes/api/render"
	"github.com/rugwirobaker/hermes/build"
	"github.com/rugwirobaker/hermes/observ"
)

// VersionHandler ...
func VersionHandler() http.HandlerFunc {
	const op = "handlers.VersionHandler"

	return func(w http.ResponseWriter, r *http.Request) {
		_, span := observ.StartSpan(r.Context(), op)
		defer span.End()

		render.JSON(w, build.Info(), http.StatusOK)
	}
}
