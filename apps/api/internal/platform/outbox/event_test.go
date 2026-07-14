package outbox

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewAndFailurePolicy(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	event, err := New("evt_1", "order", "ord_1", "ORDER_CREATED", json.RawMessage(`{"order_id":"ord_1"}`), now)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if event.Status != StatusPending {
		t.Fatalf("status = %s", event.Status)
	}
	event.MarkFailed(now, 2)
	if event.Status != StatusFailed || !event.AvailableAt.Equal(now.Add(time.Minute)) {
		t.Fatalf("unexpected first failure state: %#v", event)
	}
	event.MarkFailed(now, 2)
	if event.Status != StatusDead {
		t.Fatalf("status = %s, want DEAD", event.Status)
	}
}
