package httpapi

import (
	"net/http"
	"strings"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

func principal(r *http.Request) (application.Principal, error) {
	value := application.Principal{
		OwnerType: domain.OwnerType(strings.TrimSpace(r.Header.Get("X-Owner-Type"))),
		OwnerID:   strings.TrimSpace(r.Header.Get("X-Owner-ID")),
		ActorID:   strings.TrimSpace(r.Header.Get("X-Actor-ID")),
	}
	if value.OwnerID == "" || value.ActorID == "" {
		return application.Principal{}, domain.ErrInvalidDataset
	}
	if value.OwnerType != domain.OwnerUser && value.OwnerType != domain.OwnerOrganization {
		return application.Principal{}, domain.ErrInvalidDataset
	}
	return value, nil
}
