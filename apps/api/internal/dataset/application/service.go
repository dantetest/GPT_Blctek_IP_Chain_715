package application

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
	manifestspec "github.com/dantetest/GPT_Blctek_IP_Chain_715/packages/manifest-spec"
)

type Clock func() time.Time

type Service struct {
	repository domain.Repository
	ids        IDGenerator
	clock      Clock
}

func NewService(repository domain.Repository, ids IDGenerator, clock Clock) *Service {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &Service{repository: repository, ids: ids, clock: clock}
}

type Principal struct {
	OwnerType domain.OwnerType
	OwnerID   string
	ActorID   string
}

func (p Principal) valid() bool {
	return (p.OwnerType == domain.OwnerUser || p.OwnerType == domain.OwnerOrganization) && strings.TrimSpace(p.OwnerID) != "" && strings.TrimSpace(p.ActorID) != ""
}

type CreateDatasetCommand struct {
	Principal
	Title       string
	Slug        string
	Description string
}

type UpdateDatasetCommand struct {
	Principal
	DatasetID   string
	Title       string
	Slug        string
	Description string
	Revision    uint64
}

type CreateVersionCommand struct {
	Principal
	DatasetID    string
	VersionLabel string
}

type AttachManifestCommand struct {
	Principal
	VersionID string
	Manifest  manifestspec.Manifest
}

func (s *Service) CreateDataset(ctx context.Context, command CreateDatasetCommand) (domain.Dataset, error) {
	if !command.Principal.valid() {
		return domain.Dataset{}, domain.ErrInvalidDataset
	}
	dataset, err := domain.NewDataset(s.ids.New("dts_"), command.OwnerType, command.OwnerID, command.Title, command.Slug, command.Description, s.clock())
	if err != nil {
		return domain.Dataset{}, err
	}
	if err := s.repository.CreateDataset(ctx, dataset); err != nil {
		return domain.Dataset{}, err
	}
	return dataset, nil
}

func (s *Service) ListDatasets(ctx context.Context, principal Principal) ([]domain.Dataset, error) {
	if !principal.valid() {
		return nil, domain.ErrInvalidDataset
	}
	values, err := s.repository.ListDatasets(ctx, principal.OwnerType, principal.OwnerID)
	if err != nil {
		return nil, err
	}
	sort.Slice(values, func(i, j int) bool { return values[i].CreatedAt.After(values[j].CreatedAt) })
	return values, nil
}

func (s *Service) GetDataset(ctx context.Context, principal Principal, id string) (domain.Dataset, error) {
	dataset, err := s.repository.GetDataset(ctx, strings.TrimSpace(id))
	if err != nil {
		return domain.Dataset{}, err
	}
	if !principal.valid() || principal.OwnerType != dataset.OwnerType || principal.OwnerID != dataset.OwnerID {
		return domain.Dataset{}, domain.ErrNotFound
	}
	return dataset, nil
}

func (s *Service) UpdateDataset(ctx context.Context, command UpdateDatasetCommand) (domain.Dataset, error) {
	dataset, err := s.GetDataset(ctx, command.Principal, command.DatasetID)
	if err != nil {
		return domain.Dataset{}, err
	}
	if command.Revision == 0 || dataset.Revision != command.Revision {
		return domain.Dataset{}, domain.ErrRevisionConflict
	}
	expected := dataset.Revision
	if err := dataset.UpdateMetadata(command.Title, command.Slug, command.Description, s.clock()); err != nil {
		return domain.Dataset{}, err
	}
	if err := s.repository.UpdateDataset(ctx, dataset, expected); err != nil {
		return domain.Dataset{}, err
	}
	return dataset, nil
}

func (s *Service) CreateVersion(ctx context.Context, command CreateVersionCommand) (domain.DatasetVersion, error) {
	dataset, err := s.GetDataset(ctx, command.Principal, command.DatasetID)
	if err != nil {
		return domain.DatasetVersion{}, err
	}
	for attempt := 0; attempt < 3; attempt++ {
		number, err := s.repository.NextVersionNumber(ctx, dataset.ID)
		if err != nil {
			return domain.DatasetVersion{}, err
		}
		version, err := domain.NewDatasetVersion(s.ids.New("dsv_"), dataset.ID, number, command.VersionLabel, s.clock())
		if err != nil {
			return domain.DatasetVersion{}, err
		}
		if err := s.repository.CreateVersion(ctx, version); err != nil {
			if errors.Is(err, domain.ErrVersionNumberConflict) {
				continue
			}
			return domain.DatasetVersion{}, err
		}
		return version, nil
	}
	return domain.DatasetVersion{}, domain.ErrVersionNumberConflict
}

func (s *Service) ListVersions(ctx context.Context, principal Principal, datasetID string) ([]domain.DatasetVersion, error) {
	if _, err := s.GetDataset(ctx, principal, datasetID); err != nil {
		return nil, err
	}
	return s.repository.ListVersions(ctx, datasetID)
}

func (s *Service) GetVersion(ctx context.Context, principal Principal, versionID string) (domain.DatasetVersion, error) {
	version, err := s.repository.GetVersion(ctx, versionID)
	if err != nil {
		return domain.DatasetVersion{}, err
	}
	if _, err := s.GetDataset(ctx, principal, version.DatasetID); err != nil {
		return domain.DatasetVersion{}, err
	}
	return version, nil
}

func (s *Service) AttachManifest(ctx context.Context, command AttachManifestCommand) (domain.DatasetVersion, error) {
	if err := command.Manifest.Validate(); err != nil {
		return domain.DatasetVersion{}, errors.Join(domain.ErrManifestRequired, err)
	}
	var totalSize uint64
	for _, file := range command.Manifest.Files {
		totalSize += file.SizeBytes
	}
	root, err := domain.ParseDigest(command.Manifest.RootHash.String())
	if err != nil {
		return domain.DatasetVersion{}, err
	}
	version, err := s.GetVersion(ctx, command.Principal, command.VersionID)
	if err != nil {
		return domain.DatasetVersion{}, err
	}
	expected := version.Revision
	before := version.Status
	now := s.clock()
	if err := version.AttachManifest(domain.ManifestRef{SpecVersion: command.Manifest.ManifestVersion, Root: root, TotalFiles: uint64(len(command.Manifest.Files)), TotalSizeBytes: totalSize}, now); err != nil {
		return domain.DatasetVersion{}, err
	}
	event := domain.VersionEvent{VersionID: version.ID, From: before, To: version.Status, ActorType: "USER", ActorID: command.ActorID, Reason: "MANIFEST_ATTACHED", CreatedAt: now}
	if err := s.repository.UpdateVersion(ctx, version, expected, &event); err != nil {
		return domain.DatasetVersion{}, err
	}
	return version, nil
}
