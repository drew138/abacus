package htmx

import (
	"net/http"
	"time"
)

func (h *Handler) postLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.app.LogOut(r.Context(), sessionToken(r)); err != nil {
		h.writeError(w, r, err)
		return
	}
	h.setSessionCookie(w, "", time.Unix(0, 0))
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
