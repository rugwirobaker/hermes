package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/handlers"
	"github.com/rugwirobaker/hermes/build"
	"github.com/rugwirobaker/hermes/mock"
)

var (
	dummyPayload = &hermes.SMS{
		Payload:   "Hello",
		Recipient: "User_Phone",
	}
	dummyReport = &hermes.Report{
		ID:   "message id",
		Cost: 1,
	}
	dummyEvent = &hermes.Event{
		ID:        "fake_id",
		Status:    hermes.St(1),
		Recipient: "078xxxxxxx",
	}
	dummyCallback = &callback{
		MsgRef:     "fake_id",
		Recipient:  "078xxxxxxx",
		GatewayRef: "xxxxx",
		Status:     1,
	}
	dummyMessage = &hermes.Message{
		ID:         1,
		ProviderID: "fake_id",
		Recipient:  "078xxxxxxx",
		Payload:    "Hello",
		Cost:       1,
		Status:     hermes.St(1),
	}
)

func TestSendHander(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	sender := mock.NewMockSendService(controller)
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).Return(dummyReport, nil)

	store := mock.NewMockStore(controller)
	store.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(dummyMessage, nil)

	in := new(bytes.Buffer)

	_ = json.NewEncoder(in).Encode(dummyPayload)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	handlers.SendHandler(sender, store).ServeHTTP(w, r)
	if got, want := w.Code, http.StatusOK; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &hermes.Report{}, dummyReport
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

	got, want := &build.Build{}, build.Info()
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}
}

func TestSubscribeHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockEvent := make(chan hermes.Event)

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

	got, want := &hermes.Event{}, dummyEvent
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

	store := mock.NewMockStore(controller)
	store.EXPECT().MessageByID(gomock.Any(), gomock.Any()).Return(dummyMessage, nil)
	store.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dummyMessage, nil)

	in := new(bytes.Buffer)

	_ = json.NewEncoder(in).Encode(dummyCallback)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	handlers.DeliveryHandler(ps, store).ServeHTTP(w, r)

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
