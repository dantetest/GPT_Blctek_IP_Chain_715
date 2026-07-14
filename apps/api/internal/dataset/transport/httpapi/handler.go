package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/platform/httpx"
	manifestspec "github.com/dantetest/GPT_Blctek_IP_Chain_715/packages/manifest-spec"
)

type Handler struct{ service *application.Service }

func New(service *application.Service) *Handler { return &Handler{service: service} }

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

func principal(r *http.Request) (application.Principal, error) {
	p := application.Principal{OwnerType: domain.OwnerType(strings.TrimSpace(r.Header.Get("X-Owner-Type"))), OwnerID: strings.TrimSpace(r.Header.Get("X-Owner-ID")), ActorID: strings.TrimSpace(r.Header.Get("X-Actor-ID"))}
	if p.OwnerID == "" || p.ActorID == "" || (p.OwnerType != domain.OwnerUser && p.OwnerType != domain.OwnerOrganization) {
		return application.Principal{}, domain.ErrInvalidDataset
	}
	return p, nil
}

func decode(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, r, domain.ErrInvalidDataset)
		return false
	}
	return true
}

func write(w http.ResponseWriter, r *http.Request, status int, data any) {
	httpx.WriteJSON(w, status, httpx.Envelope{Code: "SUCCESS", Message: "success", Data: data, RequestID: httpx.RequestIDFromContext(r.Context())})
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "an internal error occurred"
	switch {
	case errors.Is(err, domain.ErrNotFound): status, code, message = http.StatusNotFound, "DATASET_NOT_FOUND", "resource not found"
	case errors.Is(err, domain.ErrRevisionConflict): status, code, message = http.StatusConflict, "DATASET_REVISION_CONFLICT", "resource revision conflict"
	case errors.Is(err, domain.ErrSlugConflict): status, code, message = http.StatusConflict, "DATASET_SLUG_CONFLICT", "dataset slug already exists"
	case errors.Is(err, domain.ErrVersionNumberConflict): status, code, message = http.StatusConflict, "DATASET_VERSION_CONFLICT", "dataset version conflict"
	case errors.Is(err, domain.ErrManifestRequired), errors.Is(err, domain.ErrInvalidDataset), errors.Is(err, domain.ErrInvalidVersion): status, code, message = http.StatusBadRequest, "INVALID_ARGUMENT", "request is invalid"
	}
	httpx.WriteJSON(w, status, httpx.Envelope{Code: code, Message: message, RequestID: httpx.RequestIDFromContext(r.Context())})
}

func (h *Handler) createDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	var body struct { Title, Slug, Description string }
	if !decode(w, r, &body) { return }
	value, err := h.service.CreateDataset(r.Context(), application.CreateDatasetCommand{Principal: p, Title: body.Title, Slug: body.Slug, Description: body.Description})
	if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusCreated, value)
}

func (h *Handler) listDatasets(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	values, err := h.service.ListDatasets(r.Context(), p); if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, values)
}

func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	value, err := h.service.GetDataset(r.Context(), p, r.PathValue("id")); if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, value)
}

func (h *Handler) updateDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	var body struct { Title, Slug, Description string; Revision uint64 }
	if !decode(w, r, &body) { return }
	value, err := h.service.UpdateDataset(r.Context(), application.UpdateDatasetCommand{Principal: p, DatasetID: r.PathValue("id"), Title: body.Title, Slug: body.Slug, Description: body.Description, Revision: body.Revision})
	if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, value)
}

func (h *Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	var body struct { VersionLabel string `json:"version_label"` }
	if !decode(w, r, &body) { return }
	value, err := h.service.CreateVersion(r.Context(), application.CreateVersionCommand{Principal: p, DatasetID: r.PathValue("id"), VersionLabel: body.VersionLabel})
	if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusCreated, value)
}

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	values, err := h.service.ListVersions(r.Context(), p, r.PathValue("id")); if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, values)
}

func (h *Handler) getVersion(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	value, err := h.service.GetVersion(r.Context(), p, r.PathValue("id")); if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, value)
}

func (h *Handler) attachManifest(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r); if err != nil { writeError(w, r, err); return }
	var manifest manifestspec.Manifest
	if !decode(w, r, &manifest) { return }
	value, err := h.service.AttachManifest(r.Context(), application.AttachManifestCommand{Principal: p, VersionID: r.PathValue("id"), Manifest: manifest})
	if err != nil { writeError(w, r, err); return }
	write(w, r, http.StatusOK, value)
}

func RevisionFromHeader(r *http.Request) (uint64, error) {
	value := strings.TrimSpace(r.Header.Get("If-Match")); value = strings.Trim(value, "\"")
	return strconv.ParseUint(value, 10, 64)
}
