package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

const datasetColumns = `id, owner_type, owner_id, title, slug, COALESCE(description, ''), status, COALESCE(default_version_id, ''), revision, created_at, updated_at`
const versionColumns = `id, dataset_id, version_number, version_label, status, manifest_spec_version, manifest_root, manifest_file_count, manifest_total_size_bytes, verification_level, revision, created_at, updated_at, published_at`

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateDataset(ctx context.Context, dataset domain.Dataset) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO datasets
		(id, owner_type, owner_id, title, slug, description, status, default_version_id, revision, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)`,
		dataset.ID, dataset.OwnerType, dataset.OwnerID, dataset.Title, dataset.Slug, dataset.Description,
		dataset.Status, dataset.DefaultVersionID, dataset.Revision, dataset.CreatedAt, dataset.UpdatedAt,
	)
	return mapWriteError(err)
}

func (r *Repository) GetDataset(ctx context.Context, id string) (domain.Dataset, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+datasetColumns+` FROM datasets WHERE id = ? AND deleted_at IS NULL`, id)
	return scanDataset(row)
}

func (r *Repository) ListDatasets(ctx context.Context, ownerType domain.OwnerType, ownerID string) ([]domain.Dataset, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+datasetColumns+` FROM datasets WHERE owner_type = ? AND owner_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`, ownerType, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := make([]domain.Dataset, 0)
	for rows.Next() {
		value, err := scanDataset(rows)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (r *Repository) UpdateDataset(ctx context.Context, dataset domain.Dataset, expectedRevision uint64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE datasets SET title = ?, slug = ?, description = ?, status = ?, default_version_id = NULLIF(?, ''), revision = ?, updated_at = ? WHERE id = ? AND revision = ? AND deleted_at IS NULL`,
		dataset.Title, dataset.Slug, dataset.Description, dataset.Status, dataset.DefaultVersionID,
		dataset.Revision, dataset.UpdatedAt, dataset.ID, expectedRevision,
	)
	if err != nil {
		return mapWriteError(err)
	}
	return requireAffected(result, domain.ErrRevisionConflict)
}

func (r *Repository) NextVersionNumber(ctx context.Context, datasetID string) (uint32, error) {
	var next uint32
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version_number), 0) + 1 FROM dataset_versions WHERE dataset_id = ?`, datasetID).Scan(&next); err != nil {
		return 0, err
	}
	return next, nil
}

func (r *Repository) CreateVersion(ctx context.Context, version domain.DatasetVersion) error {
	manifest := manifestValues(version.Manifest)
	_, err := r.db.ExecContext(ctx, `INSERT INTO dataset_versions
		(id, dataset_id, version_number, version_label, status, manifest_spec_version, manifest_root, manifest_file_count, manifest_total_size_bytes, verification_level, published_at, revision, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		version.ID, version.DatasetID, version.VersionNumber, version.VersionLabel, version.Status,
		manifest.specVersion, manifest.root, manifest.fileCount, manifest.totalSize, version.VerificationLevel,
		version.PublishedAt, version.Revision, version.CreatedAt, version.UpdatedAt,
	)
	return mapWriteError(err)
}

func (r *Repository) GetVersion(ctx context.Context, id string) (domain.DatasetVersion, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+versionColumns+` FROM dataset_versions WHERE id = ?`, id)
	return scanVersion(row)
}

func (r *Repository) ListVersions(ctx context.Context, datasetID string) ([]domain.DatasetVersion, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+versionColumns+` FROM dataset_versions WHERE dataset_id = ? ORDER BY version_number DESC`, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := make([]domain.DatasetVersion, 0)
	for rows.Next() {
		value, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (r *Repository) UpdateVersion(ctx context.Context, version domain.DatasetVersion, expectedRevision uint64, event *domain.VersionEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	manifest := manifestValues(version.Manifest)
	result, err := tx.ExecContext(ctx, `UPDATE dataset_versions SET version_label = ?, status = ?, manifest_spec_version = ?, manifest_root = ?, manifest_file_count = ?, manifest_total_size_bytes = ?, verification_level = ?, published_at = ?, revision = ?, updated_at = ? WHERE id = ? AND revision = ?`,
		version.VersionLabel, version.Status, manifest.specVersion, manifest.root, manifest.fileCount, manifest.totalSize,
		version.VerificationLevel, version.PublishedAt, version.Revision, version.UpdatedAt, version.ID, expectedRevision,
	)
	if err != nil {
		return mapWriteError(err)
	}
	if err := requireAffected(result, domain.ErrRevisionConflict); err != nil {
		return err
	}
	if event != nil {
		_, err = tx.ExecContext(ctx, `INSERT INTO dataset_version_events (dataset_version_id, from_status, to_status, actor_type, actor_id, reason, created_at) VALUES (?, NULLIF(?, ''), ?, ?, ?, ?, ?)`,
			event.VersionID, event.From, event.To, event.ActorType, event.ActorID, event.Reason, event.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

type scanner interface{ Scan(...any) error }

func scanDataset(row scanner) (domain.Dataset, error) {
	var value domain.Dataset
	if err := row.Scan(&value.ID, &value.OwnerType, &value.OwnerID, &value.Title, &value.Slug, &value.Description, &value.Status, &value.DefaultVersionID, &value.Revision, &value.CreatedAt, &value.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Dataset{}, domain.ErrNotFound
		}
		return domain.Dataset{}, err
	}
	return value, nil
}

func scanVersion(row scanner) (domain.DatasetVersion, error) {
	var value domain.DatasetVersion
	var spec sql.NullInt64
	var root []byte
	var files, size sql.NullInt64
	var published sql.NullTime
	if err := row.Scan(&value.ID, &value.DatasetID, &value.VersionNumber, &value.VersionLabel, &value.Status, &spec, &root, &files, &size, &value.VerificationLevel, &value.Revision, &value.CreatedAt, &value.UpdatedAt, &published); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.DatasetVersion{}, domain.ErrNotFound
		}
		return domain.DatasetVersion{}, err
	}
	if spec.Valid {
		if len(root) != 32 || !files.Valid || !size.Valid {
			return domain.DatasetVersion{}, fmt.Errorf("invalid persisted manifest for version %s", value.ID)
		}
		var digest domain.Digest
		copy(digest[:], root)
		value.Manifest = &domain.ManifestRef{SpecVersion: int(spec.Int64), Root: digest, TotalFiles: uint64(files.Int64), TotalSizeBytes: uint64(size.Int64)}
	}
	if published.Valid {
		value.PublishedAt = &published.Time
	}
	return value, nil
}

type nullableManifest struct {
	specVersion any
	root        any
	fileCount   any
	totalSize   any
}

func manifestValues(value *domain.ManifestRef) nullableManifest {
	if value == nil {
		return nullableManifest{}
	}
	root := make([]byte, len(value.Root))
	copy(root, value.Root[:])
	return nullableManifest{specVersion: value.SpecVersion, root: root, fileCount: value.TotalFiles, totalSize: value.TotalSizeBytes}
}

func requireAffected(result sql.Result, conflict error) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return conflict
	}
	return nil
}

func mapWriteError(err error) error {
	if err == nil {
		return nil
	}
	var mysqlErr *mysqldriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		message := mysqlErr.Message
		switch {
		case contains(message, "uq_datasets_owner_slug"):
			return domain.ErrSlugConflict
		case contains(message, "uq_dataset_versions_number"):
			return domain.ErrVersionNumberConflict
		default:
			return domain.ErrRevisionConflict
		}
	}
	return err
}

func contains(value, fragment string) bool {
	for i := 0; i+len(fragment) <= len(value); i++ {
		if value[i:i+len(fragment)] == fragment {
			return true
		}
	}
	return false
}
