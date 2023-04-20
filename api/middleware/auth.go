package middleware

import (
	"fmt"
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/render"
	"github.com/rugwirobaker/hermes/api/request"
	"github.com/rugwirobaker/hermes/observ"
)

// Authenticate  retrieves the APIKey/Token from the Bearer token in the request header and checks if we have it in our db
// the token is just a simple string, so we can just check if it exists in the db
func Authenticate(apps hermes.AppStore) Middleware {
	const op = "middleware.Authenticate"

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := observ.StartSpan(r.Context(), op)
			defer span.End()

			token := r.Header.Get("Authorization")
			if token == "" {
				err := fmt.Errorf("unauthorized: missing token")
				span.RecordError(err)
				render.HttpError(w, err)
				return
			}

			// check if the token exists in the db
			app, err := apps.FindByToken(ctx, token)
			if err != nil {
				span.RecordError(err)
				render.HttpError(w, err)
				return
			}

			ctx = request.WithApp(ctx, app)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
