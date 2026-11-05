package rest

import "net/http"

func (h *Handler) getLogin(w http.ResponseWriter, r *http.Request) {
	if _, err := h.app.Authenticate(r.Context(), sessionToken(r)); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	open, err := h.app.RegistrationOpen(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if open {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	h.render(w, r, http.StatusOK, "login", getLoginView{})
}
