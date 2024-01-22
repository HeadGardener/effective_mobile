package handlers

import (
	"net/http"
)

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("got http request", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
