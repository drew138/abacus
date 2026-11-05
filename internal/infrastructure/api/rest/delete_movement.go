package rest

import (
	"fmt"
	"net/http"

	"github.com/drew138/abacus/pkg/pesos"
)

func (h *Handler) deleteMovement(w http.ResponseWriter, r *http.Request) {
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
	if err := h.app.DeleteMovement(r.Context(), id); err != nil {
		h.writeError(w, r, err)
		return
	}
	h.showMonth(w, movement.Month(), fmt.Sprintf("Deleted %s, %s.", movement.Merchant, pesos.Format(movement.Amount)), id)
}
