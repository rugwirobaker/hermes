package render

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rugwirobaker/hermes"
)

func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}

func HttpError(w http.ResponseWriter, err error) {

	switch {
	case errors.As(err, &hermes.ErrInvalid{}):
		err, _ := err.(*hermes.ErrInvalid)
		JSON(w, err, http.StatusBadRequest)
	case err == hermes.ErrAlreadyExists:
		JSON(w, NewError(err.Error()), http.StatusConflict)
	case err == hermes.ErrNotFound:
		JSON(w, NewError(err.Error()), http.StatusNotFound)
	default:
		JSON(w, NewError("something went wrong"), http.StatusInternalServerError)
	}
}

// Flush a response down the http.ResponseWriter
func Flush(w http.ResponseWriter, f http.Flusher, v interface{}) {
	b, _ := json.Marshal(v)
	w.Write(b)
	io.WriteString(w, "\n")
	f.Flush()
}
