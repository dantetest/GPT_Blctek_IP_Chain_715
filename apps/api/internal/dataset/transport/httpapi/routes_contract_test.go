package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/infrastructure/memory"
)

func TestRegisteredDatasetRoutesDoNotReturnNotFound(t *testing.T) {
	service := application.NewService(memory.NewRepository(), application.RandomIDGenerator{}, nil)
	mux := http.NewServeMux()
	New(service).Register(mux)

	tests := []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/v1/datasets"},
		{method: http.MethodPost, path: "/api/v1/datasets"},
		{method: http.MethodGet, path: "/api/v1/datasets/dts_missing"},
		{method: http.MethodPatch, path: "/api/v1/datasets/dts_missing"},
		{method: http.MethodGet, path: "/api/v1/datasets/dts_missing/versions"},
		{method: http.MethodPost, path: "/api/v1/datasets/dts_missing/versions"},
		{method: http.MethodGet, path: "/api/v1/dataset-versions/dsv_missing"},
		{method: http.MethodPost, path: "/api/v1/dataset-versions/dsv_missing/manifest"},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)
			if res.Code == http.StatusNotFound {
				t.Fatalf("route %s %s is not registered", test.method, test.path)
			}
		})
	}
}
