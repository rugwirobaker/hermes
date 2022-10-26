package request

import (
	"context"

	"github.com/rugwirobaker/hermes"
)

type key int

const (
	requestKey key = iota
	addressKey
	idempKey
	appKey
)

// WithRequestID sets the given requestID into the context
func WithRequestID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, requestKey, id)
}

// IDFrom returns a requestID from the context or an empty
// string if not found
func IDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestKey).(string)
	return id, ok
}

func WithRemoteAddress(parent context.Context, address string) context.Context {
	return context.WithValue(parent, addressKey, address)
}

func AddressFrom(ctx context.Context) (string, bool) {
	adress, ok := ctx.Value(addressKey).(string)
	return adress, ok
}

func WithIdempotencyKey(parent context.Context, key string) context.Context {
	return context.WithValue(parent, idempKey, key)
}

func IdempotencyFrom(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(idempKey).(string)
	return key, ok
}

func WithApp(ctx context.Context, app *hermes.App) context.Context {
	return context.WithValue(ctx, appKey, app)
}

func AppFrom(ctx context.Context) (*hermes.App, bool) {
	app, ok := ctx.Value(appKey).(*hermes.App)
	return app, ok
}
