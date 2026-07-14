package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

type Repository struct {
	mu       sync.RWMutex
	datasets map[string]domain.Dataset
	versions map[string]domain.DatasetVersion
	events   []domain.VersionEvent
}

func NewRepository() *Repository {
	return &Repository{datasets: map[string]domain.Dataset{}, versions: map[string]domain.DatasetVersion{}}
}

func (r *Repository) CreateDataset(_ context.Context, dataset domain.Dataset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.datasets[dataset.ID]; exists {
		return domain.ErrRevisionConflict
	}
	for _, existing := range r.datasets {
		if existing.OwnerType == dataset.OwnerType && existing.OwnerID == dataset.OwnerID && existing.Slug == dataset.Slug {
			return domain.ErrSlugConflict
		}
	}
	r.datasets[dataset.ID] = dataset
	return nil
}

func (r *Repository) GetDataset(_ context.Context, id string) (domain.Dataset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dataset, exists := r.datasets[id]
	if !exists {
		return domain.Dataset{}, domain.ErrNotFound
	}
	return dataset, nil
}

func (r *Repository) ListDatasets(_ context.Context, ownerType domain.OwnerType, ownerID string) ([]domain.Dataset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	values := make([]domain.Dataset, 0)
	for _, dataset := range r.datasets {
		if dataset.OwnerType == ownerType && dataset.OwnerID == ownerID {
			values = append(values, dataset)
		}
	}
	sort.Slice(values, func(i, j int) bool { return values[i].CreatedAt.Before(values[j].CreatedAt) })
	return values, nil
}

func (r *Repository) UpdateDataset(_ context.Context, dataset domain.Dataset, expectedRevision uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, exists := r.datasets[dataset.ID]
	if !exists {
		return domain.ErrNotFound
	}
	if current.Revision != expectedRevision {
		return domain.ErrRevisionConflict
	}
	for _, existing := range r.datasets {
		if existing.ID != dataset.ID && existing.OwnerType == dataset.OwnerType && existing.OwnerID == dataset.OwnerID && existing.Slug == dataset.Slug {
			return domain.ErrSlugConflict
		}
	}
	r.datasets[dataset.ID] = dataset
	return nil
}

func (r *Repository) NextVersionNumber(_ context.Context, datasetID string) (uint32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var max uint32
	for _, version := range r.versions {
		if version.DatasetID == datasetID && version.VersionNumber > max {
			max = version.VersionNumber
		}
	}
	return max + 1, nil
}

func (r *Repository) CreateVersion(_ context.Context, version domain.DatasetVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.datasets[version.DatasetID]; !exists {
		return domain.ErrNotFound
	}
	if _, exists := r.versions[version.ID]; exists {
		return domain.ErrRevisionConflict
	}
	for _, existing := range r.versions {
		if existing.DatasetID == version.DatasetID && existing.VersionNumber == version.VersionNumber {
			return domain.ErrVersionNumberConflict
		}
	}
	r.versions[version.ID] = cloneVersion(version)
	return nil
}

func (r *Repository) GetVersion(_ context.Context, id string) (domain.DatasetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	version, exists := r.versions[id]
	if !exists {
		return domain.DatasetVersion{}, domain.ErrNotFound
	}
	return cloneVersion(version), nil
}

func (r *Repository) ListVersions(_ context.Context, datasetID string) ([]domain.DatasetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	values := make([]domain.DatasetVersion, 0)
	for _, version := range r.versions {
		if version.DatasetID == datasetID {
			values = append(values, cloneVersion(version))
		}
	}
	sort.Slice(values, func(i, j int) bool { return values[i].VersionNumber < values[j].VersionNumber })
	return values, nil
}

func (r *Repository) UpdateVersion(_ context.Context, version domain.DatasetVersion, expectedRevision uint64, event *domain.VersionEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, exists := r.versions[version.ID]
	if !exists {
		return domain.ErrNotFound
	}
	if current.Revision != expectedRevision {
		return domain.ErrRevisionConflict
	}
	r.versions[version.ID] = cloneVersion(version)
	if event != nil {
		r.events = append(r.events, *event)
	}
	return nil
}

func cloneVersion(value domain.DatasetVersion) domain.DatasetVersion {
	result := value
	if value.Manifest != nil {
		copy := *value.Manifest
		result.Manifest = &copy
	}
	if value.License != nil {
		copy := *value.License
		result.License = &copy
	}
	if value.Rights != nil {
		copy := *value.Rights
		result.Rights = &copy
	}
	if value.PublishedAt != nil {
		copy := *value.PublishedAt
		result.PublishedAt = &copy
	}
	return result
}
