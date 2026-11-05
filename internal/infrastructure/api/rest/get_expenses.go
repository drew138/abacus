package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getExpenses(w http.ResponseWriter, r *http.Request) {
	month := chi.URLParam(r, "month")
	if _, err := time.Parse("2006-01", month); err != nil || len(month) != 7 {
		h.writeError(w, r, errInvalidMonth)
		return
	}
	query := r.URL.Query()
	filter := expenseFilter{Category: limit(query.Get("category"), 60), Search: limit(query.Get("q"), 120)}

	report, _, err := h.loadMonthReport(r.Context(), month, filter)
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
