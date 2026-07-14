package domain

import (
	"strings"
	"time"
)

type DatasetVersion struct {
	ID                string
	DatasetID         string
	VersionNumber     uint32
	VersionLabel      string
	Status            VersionStatus
	Manifest          *ManifestRef
	License           *LicenseSnapshot
	Rights            *RightsDeclaration
	VerificationLevel VerificationLevel
	Revision          uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PublishedAt       *time.Time
}

func NewDatasetVersion(id, datasetID string, versionNumber uint32, versionLabel string, now time.Time) (DatasetVersion, error) {
	id = strings.TrimSpace(id)
	datasetID = strings.TrimSpace(datasetID)
	versionLabel = strings.TrimSpace(versionLabel)
	if id == "" || datasetID == "" || versionNumber == 0 || versionLabel == "" {
		return DatasetVersion{}, ErrInvalidVersion
	}
	now = now.UTC()
	return DatasetVersion{
		ID:                id,
		DatasetID:         datasetID,
		VersionNumber:     versionNumber,
		VersionLabel:      versionLabel,
		Status:            VersionStatusDraft,
		VerificationLevel: VerificationV0,
		Revision:          1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (v *DatasetVersion) StartScanning(now time.Time) error {
	if v.Status != VersionStatusDraft && v.Status != VersionStatusRejected {
		return ErrInvalidTransition
	}
	v.Manifest = nil
	v.Status = VersionStatusScanning
	v.touch(now)
	return nil
}

func (v *DatasetVersion) AttachManifest(manifest ManifestRef, now time.Time) error {
	if v.ContentImmutable() {
		return ErrVersionImmutable
	}
	if !v.editable() {
		return ErrInvalidTransition
	}
	if !manifest.Valid() {
		return ErrManifestRequired
	}
	copy := manifest
	v.Manifest = &copy
	v.Status = VersionStatusManifestReady
	v.touch(now)
	return nil
}

func (v *DatasetVersion) AttachLicense(license LicenseSnapshot, now time.Time) error {
	if v.ContentImmutable() {
		return ErrVersionImmutable
	}
	if !v.editable() {
		return ErrInvalidTransition
	}
	if strings.TrimSpace(license.Text) == "" || license.Hash.IsZero() {
		return ErrLicenseRequired
	}
	copy := license
	v.License = &copy
	v.touch(now)
	return nil
}

func (v *DatasetVersion) AttachRights(rights RightsDeclaration, now time.Time) error {
	if v.ContentImmutable() {
		return ErrVersionImmutable
	}
	if !v.editable() {
		return ErrInvalidTransition
	}
	if !rights.Valid() {
		return ErrRightsRequired
	}
	copy := rights
	v.Rights = &copy
	v.touch(now)
	return nil
}

func (v *DatasetVersion) SetVerificationLevel(level VerificationLevel, now time.Time) error {
	if v.ContentImmutable() {
		return ErrVersionImmutable
	}
	if !v.editable() {
		return ErrInvalidTransition
	}
	if !validVerificationLevel(level) {
		return ErrInvalidVersion
	}
	v.VerificationLevel = level
	v.touch(now)
	return nil
}

func (v *DatasetVersion) SubmitReview(now time.Time) error {
	if v.Status != VersionStatusManifestReady && v.Status != VersionStatusRejected {
		return ErrInvalidTransition
	}
	if v.Manifest == nil || !v.Manifest.Valid() {
		return ErrManifestRequired
	}
	if v.License == nil {
		return ErrLicenseRequired
	}
	if v.Rights == nil || !v.Rights.Valid() {
		return ErrRightsRequired
	}
	v.Status = VersionStatusReviewing
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Reject(now time.Time) error {
	if v.Status != VersionStatusReviewing {
		return ErrInvalidTransition
	}
	v.Status = VersionStatusRejected
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Approve(now time.Time) error {
	if v.Status != VersionStatusReviewing {
		return ErrInvalidTransition
	}
	v.Status = VersionStatusApproved
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Publish(now time.Time) error {
	if v.Status != VersionStatusApproved {
		return ErrInvalidTransition
	}
	if v.Manifest == nil || v.License == nil || v.Rights == nil {
		return ErrInvalidVersion
	}
	publishedAt := now.UTC()
	v.Status = VersionStatusPublished
	v.PublishedAt = &publishedAt
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Suspend(now time.Time) error {
	if v.Status != VersionStatusPublished {
		return ErrInvalidTransition
	}
	v.Status = VersionStatusSuspended
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Resume(now time.Time) error {
	if v.Status != VersionStatusSuspended {
		return ErrInvalidTransition
	}
	v.Status = VersionStatusPublished
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Takedown(now time.Time) error {
	if v.Status != VersionStatusPublished && v.Status != VersionStatusSuspended {
		return ErrInvalidTransition
	}
	v.Status = VersionStatusTakedown
	v.touch(now)
	return nil
}

func (v *DatasetVersion) Archive(now time.Time) error {
	switch v.Status {
	case VersionStatusRejected, VersionStatusPublished, VersionStatusSuspended, VersionStatusTakedown:
		v.Status = VersionStatusArchived
		v.touch(now)
		return nil
	default:
		return ErrInvalidTransition
	}
}

func (v DatasetVersion) editable() bool {
	switch v.Status {
	case VersionStatusDraft, VersionStatusScanning, VersionStatusManifestReady, VersionStatusRejected:
		return true
	default:
		return false
	}
}

func (v DatasetVersion) ContentImmutable() bool {
	switch v.Status {
	case VersionStatusPublished, VersionStatusSuspended, VersionStatusTakedown, VersionStatusArchived:
		return true
	default:
		return false
	}
}

func (v *DatasetVersion) touch(now time.Time) {
	v.Revision++
	v.UpdatedAt = now.UTC()
}

func validVerificationLevel(level VerificationLevel) bool {
	switch level {
	case VerificationV0, VerificationV1, VerificationV2, VerificationV3, VerificationV4:
		return true
	default:
		return false
	}
}
