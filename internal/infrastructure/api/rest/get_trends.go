package rest

import "net/http"

func (h *Handler) getTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	months, err := h.app.ListMonths(ctx)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	categories, err := h.app.ListCategories(ctx)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	totals, err := h.app.CategoryTotals(ctx)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	only := []rune(r.URL.Query().Get("only"))
	if len(only) > 60 {
		only = only[:60]
	}
	page := mapGetTrendsPage(categories, months, totals, string(only))
	if isHTMX(r) {
		h.render(w, r, http.StatusOK, "trends-swap", page)
		return
	}
	h.render(w, r, http.StatusOK, "trends-page", page)
}
