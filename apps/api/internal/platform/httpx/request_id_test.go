package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDPreservesValidCallerID(t *testing.T) {
	var observed string
	handler := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		observed = RequestIDFromContext(r.Context())
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, "caller-123")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if observed != "caller-123" || res.Header().Get(RequestIDHeader) != "caller-123" {
		t.Fatalf("request id was not propagated: observed=%q header=%q", observed, res.Header().Get(RequestIDHeader))
	}
}

func TestRequestIDReplacesUnsafeID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, "bad\nvalue")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	value := res.Header().Get(RequestIDHeader)
	if !strings.HasPrefix(value, "req_") {
		t.Fatalf("unsafe id was not replaced: %q", value)
	}
}
