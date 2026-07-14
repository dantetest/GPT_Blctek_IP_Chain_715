package memory

import (
	"context"
	"testing"
	"time"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

func TestRepositoryIsolatesOwnersAndVersions(t *testing.T) {
	repository := NewRepository()
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	dataset, err := domain.NewDataset("dts_1", domain.OwnerUser, "usr_1", "Dataset", "dataset", "", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateDataset(context.Background(), dataset); err != nil {
		t.Fatal(err)
	}
	listed, err := repository.ListDatasets(context.Background(), domain.OwnerUser, "usr_2")
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 0 {
		t.Fatalf("unexpected cross-owner datasets: %#v", listed)
	}
	version, err := domain.NewDatasetVersion("dsv_1", dataset.ID, 1, "v1", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateVersion(context.Background(), version); err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateVersion(context.Background(), version); err != domain.ErrRevisionConflict {
		t.Fatalf("duplicate id error = %v", err)
	}
}
