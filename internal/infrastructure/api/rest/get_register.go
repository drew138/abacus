package rest

import "net/http"

func (h *Handler) getRegister(w http.ResponseWriter, r *http.Request) {
	open, err := h.app.RegistrationOpen(r.Context())
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if !open {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	h.render(w, r, http.StatusOK, "register", getRegisterView{})
}
