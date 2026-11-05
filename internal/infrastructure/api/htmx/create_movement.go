package htmx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) createMovement(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(r); err != nil {
		h.writeError(w, r, err)
		return
	}
	form := mapMovementForm(r.PostForm)
	err := h.validate.Struct(form)
	var movement domain.Movement
	if err == nil {
		movement, err = h.app.CreateMovement(r.Context(), form.toDomain(0))
	}
	var failed validator.ValidationErrors
	if errors.As(err, &failed) || errors.Is(err, domain.ErrUnknownCategory) {
		h.rerenderCreateForm(w, r, form, err)
		return
	}
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.renderAfterChange(w, r, movement.Month(), toastView{Message: fmt.Sprintf("Added %s, %s, to %s.", movement.Merchant, pesos(movement.Amount), monthName(movement.Month()))})
}

func (h *Handler) rerenderCreateForm(w http.ResponseWriter, r *http.Request, form MovementForm, cause error) {
	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	draft := form.toDomain(0)
	w.Header().Set("HX-Retarget", "#modal")
	w.Header().Set("HX-Reswap", "innerHTML")
	h.render(w, r, http.StatusOK, "form", mapFormView(draft, categories, h.today(), mapFormErrors(cause)))
}
