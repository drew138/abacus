package htmx

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) getExpenses(w http.ResponseWriter, r *http.Request) {
	month := chi.URLParam(r, "month")
	if _, err := time.Parse("2006-01", month); err != nil || len(month) != 7 {
		h.writeError(w, r, errInvalidMonth)
		return
	}
	query := r.URL.Query()
	filter := domain.ExpenseFilter{Category: limit(query.Get("category"), 60), Search: limit(query.Get("q"), 120)}

	report, err := h.app.MonthReport(r.Context(), month, filter)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "expenses", mapExpensesView(report))
}

func limit(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}
