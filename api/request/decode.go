package request

import (
	"context"
	"encoding/json"
	"io"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/observ"
)

// Decode is being isolcated from the handlers so we can instrument it
// we need to close body after decoding
func Decode(ctx context.Context, body io.ReadCloser, v interface{}) error {
	const op = "request.Decode"

	defer body.Close()

	_, span := observ.StartSpan(ctx, op)
	defer span.End()

	err := json.NewDecoder(body).Decode(v)
	switch {
	case err == nil:

	case err == io.ErrUnexpectedEOF:
		err := hermes.NewErrInvalid("request body is invalid")
		span.RecordError(err)
		return err
	case err == io.EOF:
		err := hermes.NewErrInvalid("request body is empty")
		span.RecordError(err)
	case err != nil:
		err := hermes.NewErrInvalidWithErr("failed to decode request body", err)
		span.RecordError(err)
	}

	return nil
}
