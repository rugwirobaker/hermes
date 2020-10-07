package handlers

import (
	"net/http"

	"github.com/rugwirobaker/helmes"
)

// VersionHandler ...
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JSON(w, helmes.Data(), http.StatusOK)
	}
}
