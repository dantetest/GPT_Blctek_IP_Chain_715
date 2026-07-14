package domain

import (
	"strings"
	"time"
)

type Dataset struct {
	ID               string
	OwnerType        OwnerType
	OwnerID          string
	Title            string
	Slug             string
	Description      string
	Status           DatasetStatus
	DefaultVersionID string
	Revision         uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewDataset(id string, ownerType OwnerType, ownerID, title, slug, description string, now time.Time) (Dataset, error) {
	id = strings.TrimSpace(id)
	ownerID = strings.TrimSpace(ownerID)
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	if id == "" || ownerID == "" || title == "" || slug == "" || !validOwnerType(ownerType) {
		return Dataset{}, ErrInvalidDataset
	}
	now = now.UTC()
	return Dataset{
		ID:          id,
		OwnerType:   ownerType,
		OwnerID:     ownerID,
		Title:       title,
		Slug:        slug,
		Description: strings.TrimSpace(description),
		Status:      DatasetStatusActive,
		Revision:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (d *Dataset) UpdateMetadata(title, slug, description string, now time.Time) error {
	if d.Status == DatasetStatusArchived {
		return ErrInvalidDataset
	}
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	if title == "" || slug == "" {
		return ErrInvalidDataset
	}
	d.Title = title
	d.Slug = slug
	d.Description = strings.TrimSpace(description)
	d.touch(now)
	return nil
}

func (d *Dataset) SetDefaultVersion(version DatasetVersion, now time.Time) error {
	if version.DatasetID != d.ID {
		return ErrDatasetVersionMismatch
	}
	if version.Status != VersionStatusPublished {
		return ErrVersionNotPublished
	}
	d.DefaultVersionID = version.ID
	d.touch(now)
	return nil
}

func (d *Dataset) Suspend(now time.Time) error {
	if d.Status != DatasetStatusActive {
		return ErrInvalidDataset
	}
	d.Status = DatasetStatusSuspended
	d.touch(now)
	return nil
}

func (d *Dataset) Resume(now time.Time) error {
	if d.Status != DatasetStatusSuspended {
		return ErrInvalidDataset
	}
	d.Status = DatasetStatusActive
	d.touch(now)
	return nil
}

func (d *Dataset) Archive(now time.Time) error {
	if d.Status == DatasetStatusArchived {
		return ErrInvalidDataset
	}
	d.Status = DatasetStatusArchived
	d.touch(now)
	return nil
}

func (d *Dataset) touch(now time.Time) {
	d.Revision++
	d.UpdatedAt = now.UTC()
}

func validOwnerType(ownerType OwnerType) bool {
	return ownerType == OwnerUser || ownerType == OwnerOrganization
}
