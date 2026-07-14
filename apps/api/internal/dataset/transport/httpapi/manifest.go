package httpapi

import (
	"net/http"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	manifestspec "github.com/dantetest/GPT_Blctek_IP_Chain_715/packages/manifest-spec"
)

func (h *Handler) attachManifest(w http.ResponseWriter, r *http.Request) {
	p, err := principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var manifest manifestspec.Manifest
	if !decode(w, r, &manifest) {
		return
	}
	value, err := h.service.AttachManifest(r.Context(), application.AttachManifestCommand{
		Principal: p,
		VersionID: r.PathValue("id"),
		Manifest:  manifest,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	write(w, r, http.StatusOK, value)
}
