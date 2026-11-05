package rest

import "net/http"

func (h *Handler) getTrends(w http.ResponseWriter, r *http.Request) {
	months, err := h.app.ListMonths(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	categories, err := h.app.ListCategories(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	totals, err := h.app.CategoryTotals(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	view := mapTrendsView(buildTrends(categories, months, totals), limit(r.URL.Query().Get("only"), 60))
	h.renderPage(w, r, pageView{Title: "Month to month", Sidebar: mapSidebarView(months, ""), Trends: &view})
}
