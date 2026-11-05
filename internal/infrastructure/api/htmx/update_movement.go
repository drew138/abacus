package htmx

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) updateMovement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		h.writeError(w, r, errInvalidID)
		return
	}
	if err := h.parseForm(w, r); err != nil {
		h.writeError(w, r, err)
		return
	}
	form := mapMovementForm(r.PostForm)
	err = h.validate.Struct(form)
	var movement domain.Movement
	if err == nil {
		movement, err = h.app.UpdateMovement(r.Context(), form.toDomain(id), form.ApplyToMerchant)
	}
	var failed validator.ValidationErrors
	if errors.As(err, &failed) || errors.Is(err, domain.ErrUnknownCategory) {
		h.rerenderUpdateForm(w, r, id, form, err)
		return
	}
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.renderAfterChange(w, r, movement.Month(), toastView{Message: fmt.Sprintf("Saved %s.", movement.Merchant)})
}

func (h *Handler) rerenderUpdateForm(w http.ResponseWriter, r *http.Request, id int64, form MovementForm, cause error) {
	current, err := h.app.GetMovement(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	draft := form.toDomain(id)
	draft.Source = current.Source
	w.Header().Set("HX-Retarget", "#modal")
	w.Header().Set("HX-Reswap", "innerHTML")
	h.render(w, r, http.StatusOK, "form", mapFormView(draft, categories, h.today(), mapFormErrors(cause)))
}
