package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/platform/httpx"
)

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
	httpx.WriteJSON(w, status, httpx.Envelope{
		Code:      "SUCCESS",
		Message:   "success",
		Data:      data,
		RequestID: httpx.RequestIDFromContext(r.Context()),
	})
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "an internal error occurred"

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status, code, message = http.StatusNotFound, "DATASET_NOT_FOUND", "resource not found"
	case errors.Is(err, domain.ErrRevisionConflict):
		status, code, message = http.StatusConflict, "DATASET_REVISION_CONFLICT", "resource revision conflict"
	case errors.Is(err, domain.ErrSlugConflict):
		status, code, message = http.StatusConflict, "DATASET_SLUG_CONFLICT", "dataset slug already exists"
	case errors.Is(err, domain.ErrVersionNumberConflict):
		status, code, message = http.StatusConflict, "DATASET_VERSION_CONFLICT", "dataset version conflict"
	case errors.Is(err, domain.ErrManifestRequired), errors.Is(err, domain.ErrInvalidDataset), errors.Is(err, domain.ErrInvalidVersion):
		status, code, message = http.StatusBadRequest, "INVALID_ARGUMENT", "request is invalid"
	}

	httpx.WriteJSON(w, status, httpx.Envelope{
		Code:      code,
		Message:   message,
		RequestID: httpx.RequestIDFromContext(r.Context()),
	})
}
