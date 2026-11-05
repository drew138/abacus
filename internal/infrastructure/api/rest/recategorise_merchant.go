package rest

import (
	"fmt"
	"net/http"
)

func (h *Handler) recategoriseMerchant(w http.ResponseWriter, r *http.Request) {
	if err := h.parseForm(r); err != nil {
		h.writeError(w, r, err)
		return
	}
	req := RecategoriseMerchantRequest{Merchant: r.PostForm.Get("merchant"), Category: r.PostForm.Get("category"), Month: r.PostForm.Get("month")}
	if err := h.validate.Struct(req); err != nil {
		h.writeError(w, r, fmt.Errorf("%w: %w", errMalformedForm, err))
		return
	}

	changed, err := h.app.RecategoriseMerchant(r.Context(), req.Merchant, req.Category)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	plural := "s"
	if changed == 1 {
		plural = ""
	}
	h.showMonth(w, req.Month, fmt.Sprintf("%s is now %s: %d movement%s moved, and new ones will follow.", req.Merchant, req.Category, changed, plural), 0)
}
