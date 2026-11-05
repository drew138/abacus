package rest

import "net/http"

func (h *Handler) importMovements(w http.ResponseWriter, r *http.Request) {
	var req ImportRequest
	if err := h.validateRequest(w, r, &req); err != nil {
		h.writeError(w, err)
		return
	}

	movements, categories := req.toDomain()
	result, err := h.app.ImportMovements(r.Context(), movements, categories)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, mapImportResponse(result))
}
