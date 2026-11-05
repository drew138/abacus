package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func (h *Handler) createMovement(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(r); err != nil {
		h.writeError(w, r, err)
		return
	}
	req := mapCreateMovementRequest(r.PostForm)
	err := h.validate.Struct(req)
	var movement domain.Movement
	if err == nil {
		movement, err = h.app.CreateMovement(r.Context(), req.toDomain(0))
	}
	var failed validator.ValidationErrors
	if errors.As(err, &failed) || errors.Is(err, domain.ErrUnknownCategory) {
		categories, listErr := h.app.ListCategories(r.Context())
		if listErr != nil {
			h.writeError(w, r, listErr)
			return
		}
		w.Header().Set("HX-Retarget", "#modal")
		w.Header().Set("HX-Reswap", "innerHTML")
		h.render(w, r, http.StatusOK, "form", mapCreateMovementFormView(req.toDomain(0), categories, h.today(), mapCreateMovementErrors(err)))
		return
	}
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.showMonth(w, movement.Month(), fmt.Sprintf("Added %s, %s, to %s.", movement.Merchant, pesos.Format(movement.Amount), monthName(movement.Month())), 0)
}
