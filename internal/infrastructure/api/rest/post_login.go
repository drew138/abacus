package rest

import (
	"errors"
	"net/http"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) postLogin(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(r); err != nil {
		h.writeError(w, r, err)
		return
	}
	form := LoginForm{Username: r.PostForm.Get("username"), Password: r.PostForm.Get("password")}
	if err := h.validate.Struct(form); err != nil {
		h.render(w, r, http.StatusBadRequest, "login", loginView{Error: "Enter your username and password.", Username: form.Username})
		return
	}

	raw, expiresAt, err := h.app.LogIn(r.Context(), form.Username, form.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		h.render(w, r, http.StatusUnauthorized, "login", loginView{Error: "That username and password don’t match. Try again.", Username: form.Username})
		return
	}
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.setSessionCookie(w, raw, expiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
