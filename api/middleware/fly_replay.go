package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Fly.io's proxy implements a  replay(request) mechanism that allows you to
// to redirect teh request to another region, app or instance. In this middleware
// we'll take advantage of this feature to redirect the request to the primary litefs
// instance if we detect that we're not on the primary instance.
//
// to detect if we're on the primary instance we'll read litefs generated .primary file
// that contains the primary instance's IP address. If the file  doesn't exist we're either
// on the primary instance or litefs is not running on the instance.
//
// This middleware is only useful when running on fly.io
//
// For more information about the replay mechanism, please see:
// https://fly.io/docs/reference/appfile/replay/
//
// This middleware is based on the replay middleware from the chi router:

func FlyReplay(dsn string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// redict only write requests
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}
			// check if we're on the primary instance
			primary, err := isPrimary(dsn)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if primary != "" {
				log.Printf("redirecting to primary instance: %q", string(primary))
				w.Header().Set("fly-replay", "instance="+string(primary))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// isPrimary checks if we're on the primary instance by reading the .primary file
// that contains the primary instance's IP address. If the file doesn't exist we're either
// on the primary instance or litefs is not running on the instance.
func isPrimary(dsn string) (string, error) {
	primaryFilename := filepath.Join(filepath.Dir(dsn), ".primary")

	primary, err := os.ReadFile(primaryFilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return string(primary), nil
}
