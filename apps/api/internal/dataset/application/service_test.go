package application

import (
	"context"
	"testing"
	"time"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/infrastructure/memory"
)

type fixedIDs struct{ values []string }

func (f *fixedIDs) New(string) string {
	value := f.values[0]
	f.values = f.values[1:]
	return value
}

func TestServiceEnforcesOwnershipAndOptimisticRevision(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	service := NewService(memory.NewRepository(), &fixedIDs{values: []string{"dts_test"}}, func() time.Time { return now })
	owner := Principal{OwnerType: domain.OwnerUser, OwnerID: "usr_owner", ActorID: "usr_owner"}
	dataset, err := service.CreateDataset(context.Background(), CreateDatasetCommand{Principal: owner, Title: "Dataset", Slug: "dataset"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetDataset(context.Background(), Principal{OwnerType: domain.OwnerUser, OwnerID: "usr_other", ActorID: "usr_other"}, dataset.ID); err != domain.ErrNotFound {
		t.Fatalf("cross-owner read error = %v", err)
	}
	if _, err := service.UpdateDataset(context.Background(), UpdateDatasetCommand{Principal: owner, DatasetID: dataset.ID, Title: "Updated", Slug: "updated", Revision: 99}); err != domain.ErrRevisionConflict {
		t.Fatalf("revision error = %v", err)
	}
}
