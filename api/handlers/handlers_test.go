package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
