package domain

import "context"

type Repository interface {
	CreateDataset(ctx context.Context, dataset Dataset) error
	GetDataset(ctx context.Context, id string) (Dataset, error)
	UpdateDataset(ctx context.Context, dataset Dataset, expectedRevision uint64) error

	NextVersionNumber(ctx context.Context, datasetID string) (uint32, error)
	CreateVersion(ctx context.Context, version DatasetVersion) error
	GetVersion(ctx context.Context, id string) (DatasetVersion, error)
	UpdateVersion(ctx context.Context, version DatasetVersion, expectedRevision uint64) error
}
