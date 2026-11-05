package rest

import "net/http"

func (h *Handler) importMovements(w http.ResponseWriter, r *http.Request) {
	var req ImportRequest
	if err := h.validateRequest(r, &req); err != nil {
		h.writeAPIError(w, err)
		return
	}

	movements, categories := req.toDomain()
	added, updated, kept, err := h.app.ImportMovements(r.Context(), movements, categories)
	if err != nil {
		h.writeAPIError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, mapImportResponse(added, updated, kept))
}
