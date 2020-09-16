package main

import (
	"net/http"

	"github.com/quarksgroup/sms-client/sms"
	"github.com/quarksgroup/sms-client/sms/driver/fdi"
	"github.com/quarksgroup/sms-client/sms/transport/oauth2"
)

func provideClient() *sms.Client {
	client := fdi.NewDefault()
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Scheme: oauth2.SchemeBearer,
			Source: oauth2.ContextTokenSource(),
			Base:   http.DefaultTransport,
		},
	}
	return client
}
