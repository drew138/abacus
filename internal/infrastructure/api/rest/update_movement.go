package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) updateMovement(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if err := h.parseForm(r); err != nil {
		h.writeError(w, r, err)
		return
	}
	req := mapUpdateMovementRequest(r.PostForm)
	err = h.validate.Struct(req)
	var movement domain.Movement
	if err == nil {
		movement, err = h.app.UpdateMovement(r.Context(), req.toDomain(id), req.ApplyToMerchant)
	}
	var failed validator.ValidationErrors
	if errors.As(err, &failed) || errors.Is(err, domain.ErrUnknownCategory) {
		current, getErr := h.app.GetMovement(r.Context(), id)
		if getErr != nil {
			h.writeError(w, r, getErr)
			return
		}
		categories, listErr := h.app.ListCategories(r.Context())
		if listErr != nil {
			h.writeError(w, r, listErr)
			return
		}
		draft := req.toDomain(id)
		draft.Source = current.Source
		w.Header().Set("HX-Retarget", "#modal")
		w.Header().Set("HX-Reswap", "innerHTML")
		h.render(w, r, http.StatusOK, "form", mapUpdateMovementFormView(draft, categories, h.today(), mapUpdateMovementErrors(err)))
		return
	}
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.showMonth(w, movement.Month(), fmt.Sprintf("Saved %s.", movement.Merchant), 0)
}
