package hermes

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/quarksgroup/sms-client/sms"
	"github.com/rugwirobaker/hermes/observ"
)

// SMS ...
type SMS struct {
	Sender    string `json:"sender"`
	Payload   string `json:"payload"`
	Recipient string `json:"recipient"`
}

// Report message queueing status
type Report struct {
	ID   string `json:"id"`
	Cost int64  `json:"cost"`
}

type Callback struct {
	MsgRef     string `json:"msgRef"`
	Recipient  string `json:"recipient"`
	GatewayRef string `json:"gatewayRef"`
	Status     int    `json:"status"`
}

// SendService is a front to the sending service
type SendService interface {
	// Send an sms message and return it's
	Send(context.Context, *SMS) (*Report, error)
}

type service struct {
	callback string
	client   *sms.Client
	token    *sms.Token
}

// NewSendService instance of service
func NewSendService(cli *sms.Client, id, secret, sender, callback string) (SendService, error) {
	token, _, err := cli.Auth.Login(context.Background(), id, secret)
	if err != nil {
		return nil, err
	}
	return &service{
		callback: callback,
		client:   cli,
		token:    token,
	}, nil
}

func (s *service) Send(ctx context.Context, message *SMS) (*Report, error) {
	const op = "service.Send"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	token, _, err := s.client.Auth.Refresh(ctx, s.token, false)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	s.token = token
	ctx = context.WithValue(ctx, sms.TokenKey{}, &sms.Token{
		Token:   token.Token,
		Refresh: token.Refresh,
	})

	log.Printf("sender: %s", message.Sender)

	in := sms.Message{
		ID:         uuid.New().String(),
		Body:       message.Payload,
		Recipients: []string{message.Recipient},
		Sender:     message.Sender,
		Report:     s.callback,
	}

	report, _, err := s.client.Message.Send(ctx, in)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return convertReport(report), nil
}

func convertReport(report *sms.Report) *Report {
	return &Report{
		ID:   report.ID,
		Cost: report.Cost,
	}
}
