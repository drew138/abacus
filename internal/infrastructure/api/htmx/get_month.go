package htmx

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) getMonth(w http.ResponseWriter, r *http.Request) {
	month := chi.URLParam(r, "month")
	if _, err := time.Parse("2006-01", month); err != nil || len(month) != 7 {
		h.writeError(w, r, errInvalidMonth)
		return
	}

	page, err := h.monthPage(r.Context(), month, domain.ExpenseFilter{})
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.renderPage(w, r, page)
}
