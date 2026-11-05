package htmx

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) deleteMovement(w http.ResponseWriter, r *http.Request) {
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
	if err := h.app.DeleteMovement(r.Context(), id); err != nil {
		h.writeError(w, r, err)
		return
	}
	h.renderAfterChange(w, r, movement.Month(), toastView{Message: fmt.Sprintf("Deleted %s, %s.", movement.Merchant, pesos(movement.Amount)), UndoID: id})
}
