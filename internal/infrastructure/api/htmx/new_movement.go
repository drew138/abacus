package htmx

import (
	"net/http"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) newMovement(w http.ResponseWriter, r *http.Request) {
	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	today := h.today()
	date := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	if month, err := time.Parse("2006-01", r.URL.Query().Get("month")); err == nil && month.Format("2006-01") != date.Format("2006-01") && month.Before(date) {
		date = month
	}
	movement := domain.Movement{Date: date, Method: domain.MethodCredit, Category: "Restaurants and cafés"}
	h.render(w, r, http.StatusOK, "form", mapFormView(movement, categories, today, nil))
}
