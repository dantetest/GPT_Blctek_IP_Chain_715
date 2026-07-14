package httpapi

import (
	"net/http"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
)

func (h *Handler) createDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var body struct {
		Title       string
		Slug        string
		Description string
	}
	if !decode(w, r, &body) {
		return
	}
	value, err := h.service.CreateDataset(r.Context(), application.CreateDatasetCommand{
		Principal:   p,
		Title:       body.Title,
		Slug:        body.Slug,
		Description: body.Description,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusCreated, value)
}

func (h *Handler) listDatasets(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	values, err := h.service.ListDatasets(r.Context(), p)
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, values)
}

func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	value, err := h.service.GetDataset(r.Context(), p, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, value)
}

func (h *Handler) updateDataset(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var body struct {
		Title       string
		Slug        string
		Description string
		Revision    uint64
	}
	if !decode(w, r, &body) {
		return
	}
	value, err := h.service.UpdateDataset(r.Context(), application.UpdateDatasetCommand{
		Principal:   p,
		DatasetID:   r.PathValue("id"),
		Title:       body.Title,
		Slug:        body.Slug,
		Description: body.Description,
		Revision:    body.Revision,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, value)
}
