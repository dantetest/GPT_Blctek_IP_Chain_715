package mysqlrepo

import (
	"errors"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

func TestMapWriteErrorClassifiesUniqueConstraints(t *testing.T) {
	tests := []struct {
		message string
		want    error
	}{
		{message: "uq_datasets_owner_slug", want: domain.ErrSlugConflict},
		{message: "uq_dataset_versions_number", want: domain.ErrVersionNumberConflict},
		{message: "PRIMARY", want: domain.ErrRevisionConflict},
	}
	for _, test := range tests {
		err := mapWriteError(&mysqldriver.MySQLError{Number: 1062, Message: test.message})
		if !errors.Is(err, test.want) {
			t.Fatalf("mapWriteError() = %v, want %v", err, test.want)
		}
	}
}
