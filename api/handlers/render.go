package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSON responds with json
func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}

func HttpError(w http.ResponseWriter, err error, status int) {
	JSON(w, NewError(err.Error()), status)
}

// Flush a response down the http.ResponseWriter
func Flush(w http.ResponseWriter, f http.Flusher, v interface{}) {
	b, _ := json.Marshal(v)
	w.Write(b)
	io.WriteString(w, "\n")
	f.Flush()
}
