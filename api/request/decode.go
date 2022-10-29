package request

import (
	"context"
	"encoding/json"
	"io"

	"github.com/rugwirobaker/hermes/observ"
)

// Decode is being isolcated from the handlers so we can instrument it
// we need to close body after decoding
func Decode(ctx context.Context, body io.ReadCloser, v interface{}) error {
	const op = "request.Decode"

	defer body.Close()

	_, span := observ.StartSpan(ctx, op)
	defer span.End()

	if err := json.NewDecoder(body).Decode(v); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
