package handlers

import (
	"net/http"

	"github.com/rugwirobaker/hermes"
)

// VersionHandler ...
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JSON(w, hermes.Data(), http.StatusOK)
	}
}
