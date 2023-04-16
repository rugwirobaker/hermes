package hermes

import (
	"context"
	"log"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/pindo"
)

// SMS ...
type SMS struct {
	Sender    string `json:"sender"`
	Payload   string `json:"payload"`
	Recipient string `json:"recipient"`
}

// Report message queueing status
type Report struct {
	Count  int64   `json:"count"`
	ID     int     `json:"id"`
	Cost   float64 `json:"cost"`
	Status string  `json:"status"`
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
	sender string
	client *pindo.Client
}

// NewSendService instance of service
func NewSendService(cli *pindo.Client, sender string) (SendService, error) {
	return &service{
		sender: sender,
		client: cli,
	}, nil
}

func (s *service) Send(ctx context.Context, message *SMS) (*Report, error) {
	const op = "service.Send"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	sendRequest := &pindo.SendRequest{
		Sender: message.Sender,
		To:     message.Recipient,
		Text:   message.Payload,
	}

	if sendRequest.Sender == "" {
		sendRequest.Sender = s.sender
	}

	log.Println("sending message", sendRequest)

	res, err := s.client.Send(ctx, sendRequest)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return convertResponse(res), nil
}

func convertResponse(res *pindo.SendResponse) *Report {
	return &Report{
		Count:  res.ItemCount,
		ID:     res.SmsID,
		Status: res.Status,
		Cost:   res.TotalCost,
	}
}
