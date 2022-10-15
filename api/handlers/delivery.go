package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rugwirobaker/hermes"
)

// DeliveryHandler handles delivery callback reception
func DeliveryHandler(events hermes.Pubsub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Body == nil {
			log.Printf("--> %s %s %s\n", r.Method, r.URL.String(), r.Header.Get("Agent"))
		} else {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			r.Body = ioutil.NopCloser(buf)
			log.Printf("--> %s %s %s %s\n", r.Method, r.URL.String(), r.Header.Get("User-Agent"), buf.String())
		}

		var in map[string]interface{}

		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Printf("failed to serialize request")
			JSON(w, NewError(err.Error()), 500)
			return
		}
		// events.Publish(r.Context(), convertEvent(in))

		JSON(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}
