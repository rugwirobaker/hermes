package sam

import (
	"context"
)

// SMS ...
type SMS struct {
	Payload   string `json:"payload"`
	Recipient string `json:"recipient"`
}

// Service defines the capabilties of sam
type Service interface {
	// Send an sms message and return it's
	Send(context.Context, SMS) (string, int, error)

	//Version returns sam's current running version
	Version(context.Context) (string, error)
}

type service struct{}

// New instance of service
func New() Service {
	return &service{}
}

func (s *service) Send(ctx context.Context, sms SMS) (string, int, error) {
	return "ok", 0, nil
}

func (s *service) Version(ctx context.Context) (string, error) {
	return "v0.1.0", nil
}
