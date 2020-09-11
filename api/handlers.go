package api

import (
	"fmt"
	"net/http"

	"github.com/rugwirobaker/sam"
)

// SMSHandler ...
func SMSHandler(svc sam.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, cost, err := svc.Send(r.Context(), sam.SMS{})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte(fmt.Sprintf("%s:%d", id, cost)))
	}
}

// VersionHandler ...
func VersionHandler(svc sam.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ver, err := svc.Version(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte(ver))
	}
}
