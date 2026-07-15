package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

func TestRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		t.Skip("MYSQL_TEST_DSN is not set")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	repository := New(db)
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	dataset, err := domain.NewDataset("dts_mysql_integration", domain.OwnerUser, "usr_mysql_integration", "MySQL integration", "mysql-integration", "", now)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = db.ExecContext(ctx, "DELETE FROM dataset_version_events WHERE dataset_version_id = ?", "dsv_mysql_integration")
	_, _ = db.ExecContext(ctx, "DELETE FROM dataset_versions WHERE id = ?", "dsv_mysql_integration")
	_, _ = db.ExecContext(ctx, "DELETE FROM datasets WHERE id = ?", dataset.ID)
	defer func() {
		_, _ = db.Exec("DELETE FROM dataset_version_events WHERE dataset_version_id = ?", "dsv_mysql_integration")
		_, _ = db.Exec("DELETE FROM dataset_versions WHERE id = ?", "dsv_mysql_integration")
		_, _ = db.Exec("DELETE FROM datasets WHERE id = ?", dataset.ID)
	}()

	if err := repository.CreateDataset(ctx, dataset); err != nil {
		t.Fatal(err)
	}
	loaded, err := repository.GetDataset(ctx, dataset.ID)
	if err != nil || loaded.Slug != dataset.Slug {
		t.Fatalf("GetDataset() = %#v, %v", loaded, err)
	}
	loaded.Title = "Updated"
	loaded.Revision++
	loaded.UpdatedAt = now.Add(time.Minute)
	if err := repository.UpdateDataset(ctx, loaded, dataset.Revision); err != nil {
		t.Fatal(err)
	}
	if err := repository.UpdateDataset(ctx, loaded, dataset.Revision); !errors.Is(err, domain.ErrRevisionConflict) {
		t.Fatalf("stale update error = %v", err)
	}

	version, err := domain.NewDatasetVersion("dsv_mysql_integration", dataset.ID, 1, "v1", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateVersion(ctx, version); err != nil {
		t.Fatal(err)
	}
	root, err := domain.ParseDigest("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	expectedRevision := version.Revision
	if err := version.AttachManifest(domain.ManifestRef{SpecVersion: 1, Root: root, TotalFiles: 2, TotalSizeBytes: 42}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	event := &domain.VersionEvent{VersionID: version.ID, From: domain.VersionStatusDraft, To: version.Status, ActorType: "USER", ActorID: dataset.OwnerID, Reason: "INTEGRATION_TEST", CreatedAt: now.Add(time.Minute)}
	if err := repository.UpdateVersion(ctx, version, expectedRevision, event); err != nil {
		t.Fatal(err)
	}
	persisted, err := repository.GetVersion(ctx, version.ID)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.Manifest == nil || persisted.Manifest.Root != root || persisted.Status != domain.VersionStatusManifestReady {
		t.Fatalf("persisted version = %#v", persisted)
	}
	var eventCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM dataset_version_events WHERE dataset_version_id = ?", version.ID).Scan(&eventCount); err != nil {
		t.Fatal(err)
	}
	if eventCount != 1 {
		t.Fatalf("event count = %d", eventCount)
	}
}
