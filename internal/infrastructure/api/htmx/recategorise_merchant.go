package htmx

import (
	"fmt"
	"net/http"
)

func (h *Handler) recategoriseMerchant(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(w, r); err != nil {
		h.writeError(w, r, err)
		return
	}
	form := RecategoriseForm{Merchant: r.PostForm.Get("merchant"), Category: r.PostForm.Get("category"), Month: r.PostForm.Get("month")}
	if err := h.validate.Struct(form); err != nil {
		h.writeError(w, r, fmt.Errorf("%w: %w", errMalformedForm, err))
		return
	}

	changed, err := h.app.RecategoriseMerchant(r.Context(), form.Merchant, form.Category)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	plural := "s"
	if changed == 1 {
		plural = ""
	}
	h.renderAfterChange(w, r, form.Month, toastView{Message: fmt.Sprintf("%s is now %s: %d movement%s moved, and new ones will follow.", form.Merchant, form.Category, changed, plural)})
}
