package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/infrastructure/memory"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/platform/httpx"
)

func newTestHandler() http.Handler {
	repository := memory.NewRepository()
	service := application.NewService(repository, application.RandomIDGenerator{}, func() time.Time {
		return time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	})
	mux := http.NewServeMux()
	New(service).Register(mux)
	return httpx.RequestID(mux)
}

func perform(handler http.Handler, method, path, body, owner string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	if owner != "" {
		request.Header.Set("X-Owner-Type", "USER")
		request.Header.Set("X-Owner-ID", owner)
		request.Header.Set("X-Actor-ID", owner)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestDatasetOwnerIsolation(t *testing.T) {
	handler := newTestHandler()
	created := perform(handler, http.MethodPost, "/api/v1/datasets", `{"Title":"Images","Slug":"images","Description":"training images"}`, "usr_a")
	if created.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
	}
	var payload struct {
		Data struct{ ID string } `json:"data"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	foreign := perform(handler, http.MethodGet, "/api/v1/datasets/"+payload.Data.ID, "", "usr_b")
	if foreign.Code != http.StatusNotFound {
		t.Fatalf("foreign status=%d body=%s", foreign.Code, foreign.Body.String())
	}
}

func TestDatasetRequiresPrincipal(t *testing.T) {
	response := perform(newTestHandler(), http.MethodGet, "/api/v1/datasets", "", "")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}
