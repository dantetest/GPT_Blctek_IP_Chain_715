package domain

import (
	"context"
	"errors"
)

var (
	ErrNotFound              = errors.New("dataset resource not found")
	ErrRevisionConflict      = errors.New("dataset revision conflict")
	ErrSlugConflict          = errors.New("dataset slug already exists")
	ErrVersionNumberConflict = errors.New("dataset version number conflict")
)

type Repository interface {
	CreateDataset(ctx context.Context, dataset Dataset) error
	GetDataset(ctx context.Context, id string) (Dataset, error)
	ListDatasets(ctx context.Context, ownerType OwnerType, ownerID string) ([]Dataset, error)
	UpdateDataset(ctx context.Context, dataset Dataset, expectedRevision uint64) error

	NextVersionNumber(ctx context.Context, datasetID string) (uint32, error)
	CreateVersion(ctx context.Context, version DatasetVersion) error
	GetVersion(ctx context.Context, id string) (DatasetVersion, error)
	ListVersions(ctx context.Context, datasetID string) ([]DatasetVersion, error)
	UpdateVersion(ctx context.Context, version DatasetVersion, expectedRevision uint64, event *VersionEvent) error
}
