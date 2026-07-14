package problem

import (
	"net/http"
	"testing"
)

func TestConflict(t *testing.T) {
	got := Conflict("ORDER_INVALID_STATE", "invalid transition")
	if got.Status != http.StatusConflict || got.Code != "ORDER_INVALID_STATE" {
		t.Fatalf("unexpected problem: %#v", got)
	}
}
