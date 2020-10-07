package helmes_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/sms-client/sms"
	"github.com/rugwirobaker/helmes"
	"github.com/rugwirobaker/helmes/mock/mocksmc"
)

var noContext = context.Background()

func TestSend(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockReport := &sms.Report{
		ID:   "fake_id",
		Cost: 1,
	}
	mockToken := &sms.Token{}
	mockSMS := &helmes.SMS{}

	mockSendService := mocksmc.NewMockSendService(controller)
	mockSendService.EXPECT().Send(gomock.Any(), gomock.Any()).Return(mockReport, nil, nil)

	mockAuthService := mocksmc.NewMockAuthService(controller)
	mockAuthService.EXPECT().Refresh(gomock.Any(), gomock.Any(), false).Return(mockToken, nil, nil)
	mockAuthService.EXPECT().Login(gomock.Any(), "id", "secret").Return(mockToken, nil, nil)

	client := new(sms.Client)
	client.Message = mockSendService
	client.Auth = mockAuthService

	service, _ := helmes.NewSendService(client, "id", "secret", "sender", "callback")

	want := &helmes.Report{
		ID:   "fake_id",
		Cost: 1,
	}

	got, err := service.Send(noContext, mockSMS)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}

}
