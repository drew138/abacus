package htmx

import "net/http"

func (h *Handler) getTrends(w http.ResponseWriter, r *http.Request) {
	months, err := h.app.ListMonths(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	trends, err := h.app.Trends(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	view := mapTrendsView(trends, limit(r.URL.Query().Get("only"), 60))
	h.renderPage(w, r, pageView{Title: "Month to month", Sidebar: mapSidebarView(months, ""), Trends: &view})
}
