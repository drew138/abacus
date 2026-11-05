package htmx

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

func (h *Handler) postRegister(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(w, r); err != nil {
		h.writeError(w, r, err)
		return
	}
	form := RegisterForm{Username: r.PostForm.Get("username"), Password: r.PostForm.Get("password"), Confirm: r.PostForm.Get("confirm")}
	if err := h.validate.Struct(form); err != nil {
		h.render(w, r, http.StatusBadRequest, "register", loginView{Error: mapRegisterError(err), Username: form.Username})
		return
	}

	_, err := h.app.Register(r.Context(), form.Username, form.Password)
	switch {
	case errors.Is(err, domain.ErrRegistrationClosed):
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	case errors.Is(err, domain.ErrPasswordTooShort):
		h.render(w, r, http.StatusBadRequest, "register", loginView{Error: mapRegisterError(err), Username: form.Username})
		return
	case err != nil:
		h.writeError(w, r, err)
		return
	}

	raw, session, err := h.app.LogIn(r.Context(), form.Username, form.Password)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.setSessionCookie(w, raw, session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func mapRegisterError(err error) string {
	var failed validator.ValidationErrors
	if errors.As(err, &failed) {
		for _, fe := range failed {
			switch fe.Field() {
			case "username":
				return "Choose a username of up to 60 characters."
			case "confirm":
				return "The two passwords differ. Type the same password twice."
			}
		}
	}
	return "Use a password of 12 to 72 characters."
}
