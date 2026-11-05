package rest

import "net/http"

func (h *Handler) getRegisterPage(w http.ResponseWriter, r *http.Request) {
	open, err := h.app.CheckRegistration(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if !open {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	h.render(w, r, http.StatusOK, "register", getRegisterPageView{})
}
