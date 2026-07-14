package outbox

import (
	"encoding/json"
	"errors"
	"time"
)

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusProcessed  Status = "PROCESSED"
	StatusFailed     Status = "FAILED"
	StatusDead       Status = "DEAD"
)

var ErrInvalidEvent = errors.New("outbox event is invalid")

type Event struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       json.RawMessage
	Status        Status
	Attempts      int
	AvailableAt   time.Time
}

func New(id, aggregateType, aggregateID, eventType string, payload json.RawMessage, now time.Time) (Event, error) {
	if id == "" || aggregateType == "" || aggregateID == "" || eventType == "" || !json.Valid(payload) {
		return Event{}, ErrInvalidEvent
	}
	return Event{
		ID:            id,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       append(json.RawMessage(nil), payload...),
		Status:        StatusPending,
		AvailableAt:   now.UTC(),
	}, nil
}

func (e *Event) MarkFailed(now time.Time, maxAttempts int) {
	e.Attempts++
	if e.Attempts >= maxAttempts {
		e.Status = StatusDead
		return
	}
	e.Status = StatusFailed
	e.AvailableAt = now.UTC().Add(backoff(e.Attempts))
}

func backoff(attempt int) time.Duration {
	delays := []time.Duration{time.Minute, 5 * time.Minute, 15 * time.Minute, time.Hour, 6 * time.Hour, 24 * time.Hour}
	if attempt <= 0 {
		return delays[0]
	}
	if attempt > len(delays) {
		return delays[len(delays)-1]
	}
	return delays[attempt-1]
}
