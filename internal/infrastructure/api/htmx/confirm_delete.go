package htmx

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) confirmDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		h.writeError(w, r, errInvalidID)
		return
	}
	movement, err := h.app.GetMovement(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "row", mapRowView(movement, true))
}
