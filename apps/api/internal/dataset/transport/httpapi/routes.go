package httpapi

import (
	"net/http"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
)

type Handler struct {
	service *application.Service
}

func New(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/datasets", h.listDatasets)
	mux.HandleFunc("POST /api/v1/datasets", h.createDataset)
	mux.HandleFunc("GET /api/v1/datasets/{id}", h.getDataset)
	mux.HandleFunc("PATCH /api/v1/datasets/{id}", h.updateDataset)
	mux.HandleFunc("GET /api/v1/datasets/{id}/versions", h.listVersions)
	mux.HandleFunc("POST /api/v1/datasets/{id}/versions", h.createVersion)
	mux.HandleFunc("GET /api/v1/dataset-versions/{id}", h.getVersion)
	mux.HandleFunc("POST /api/v1/dataset-versions/{id}/manifest", h.attachManifest)
}
