package htmx

import "net/http"

func (h *Handler) getLogin(w http.ResponseWriter, r *http.Request) {
	if _, err := h.app.Authenticate(r.Context(), sessionToken(r)); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.render(w, r, http.StatusOK, "login", loginView{})
}
