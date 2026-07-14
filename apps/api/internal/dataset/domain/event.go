package domain

import "time"

type VersionEvent struct {
	VersionID string
	From      VersionStatus
	To        VersionStatus
	ActorType string
	ActorID   string
	Reason    string
	CreatedAt time.Time
}
