package htmx

import "net/http"

func (h *Handler) getIndex(w http.ResponseWriter, r *http.Request) {
	months, err := h.app.ListMonths(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if len(months) == 0 {
		http.Redirect(w, r, "/trends", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/months/"+months[len(months)-1].Month, http.StatusSeeOther)
}
