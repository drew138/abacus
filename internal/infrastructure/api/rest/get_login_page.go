package rest

import "net/http"

func (h *Handler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	if _, err := h.app.AuthenticateUser(r.Context(), sessionToken(r)); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	open, err := h.app.CheckRegistration(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if open {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	h.render(w, r, http.StatusOK, "login", getLoginPageView{})
}
