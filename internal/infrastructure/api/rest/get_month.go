package rest

import "net/http"

func (h *Handler) getMonth(w http.ResponseWriter, r *http.Request) {
	month, err := parseMonth(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
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
	movements, err := h.app.ListMovements(ctx, month)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	page := mapGetMonthPage(month, categories, months, totals, movements)
	if isHTMX(r) {
		h.render(w, r, http.StatusOK, "month-swap", page)
		return
	}
	h.render(w, r, http.StatusOK, "month-page", page)
}
