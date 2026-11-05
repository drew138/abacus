package rest

import "net/http"

func (h *Handler) getExpenses(w http.ResponseWriter, r *http.Request) {
	month, err := parseMonth(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	query := getExpensesQuery{Category: limitGetExpenses(r.URL.Query().Get("category"), 60), Search: limitGetExpenses(r.URL.Query().Get("q"), 120)}

	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	movements, err := h.app.ListMovements(r.Context(), month)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "expenses", mapGetExpensesView(month, query, categories, movements))
}

func limitGetExpenses(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}
