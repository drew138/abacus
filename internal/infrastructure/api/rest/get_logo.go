package rest

import "net/http"

func (h *Handler) getLogo(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/static/logo.svg"
	h.getStaticFile(w, r)
}
