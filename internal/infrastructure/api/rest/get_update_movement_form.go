package rest

import "net/http"

func (h *Handler) getUpdateMovementForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	movement, err := h.app.GetMovement(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "form", mapGetUpdateMovementFormView(movement, categories, h.today(), nil))
}
