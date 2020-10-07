package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/rugwirobaker/helmes"
	"github.com/rugwirobaker/helmes/api/handlers"
	"github.com/rugwirobaker/helmes/mock"
)

var (
	dummyMessage = &helmes.SMS{
		Payload:   "Hello",
		Recipient: "User_Phone",
	}
	dummyReport = &helmes.Report{
		ID:   "message id",
		Cost: 1,
	}
	dummyEvent = &helmes.Event{
		ID:        "fake_id",
		Status:    helmes.St(1),
		Recipient: "078xxxxxxx",
	}
	dummyCallback = &callback{
		MsgRef:     "fake_id",
		Recipient:  "078xxxxxxx",
		GatewayRef: "xxxxx",
		Status:     1,
	}
)

func TestSendHander(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	sender := mock.NewMockSendService(controller)
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).Return(dummyReport, nil)

	in := new(bytes.Buffer)

	_ = json.NewEncoder(in).Encode(dummyMessage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	handlers.SendHandler(sender).ServeHTTP(w, r)
	if got, want := w.Code, http.StatusOK; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &helmes.Report{}, dummyReport
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}

}

func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/healthz", nil)

	handlers.HealthHandler().ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestVersionHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/version", nil)

	handlers.VersionHandler().ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &helmes.Build{}, helmes.Data()
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}
}

func TestSubscribeHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockEvent := make(chan helmes.Event)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		mockEvent <- *dummyEvent
		wg.Done()
	}()

	ps := mock.NewMockPubsub(controller)
	ps.EXPECT().Subscribe(gomock.Any(), "fake_id").Return(mockEvent, nil)
	ps.EXPECT().Done(gomock.Any(), "fake_id").Return(nil)

	c := new(chi.Context)
	c.URLParams.Add("id", "fake_id")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		context.WithValue(context.Background(), chi.RouteCtxKey, c),
	)

	handlers.SubscribeHandler(ps).ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &helmes.Event{}, dummyEvent
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}
}

func TestDeliveryHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ps := mock.NewMockPubsub(controller)
	ps.EXPECT().Publish(gomock.Any(), gomock.Any())

	in := new(bytes.Buffer)

	_ = json.NewEncoder(in).Encode(dummyCallback)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	handlers.DeliveryHandler(ps).ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

type callback struct {
	MsgRef     string `json:"msgRef"`
	Recipient  string `json:"recipient"`
	GatewayRef string `json:"gatewayRef"`
	Status     int    `json:"status"`
}
