package httpapi

import (
	"net/http"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
)

func (h *Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var body struct {
		VersionLabel string `json:"version_label"`
	}
	if !decode(w, r, &body) {
		return
	}
	value, err := h.service.CreateVersion(r.Context(), application.CreateVersionCommand{
		Principal:    p,
		DatasetID:    r.PathValue("id"),
		VersionLabel: body.VersionLabel,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusCreated, value)
}

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	values, err := h.service.ListVersions(r.Context(), p, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, values)
}

func (h *Handler) getVersion(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	value, err := h.service.GetVersion(r.Context(), p, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, value)
}
