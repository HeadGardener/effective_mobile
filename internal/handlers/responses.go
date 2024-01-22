package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/HeadGardener/effective_mobile/internal/services"
)

type response struct {
	Msg   string `json:"Msg"`
	Error string `json:"Error"`
}

func (h *Handler) newErrResponse(w http.ResponseWriter, code int, msg string, err error) {
	h.log.Error(msg, "error", err.Error())

	if !errIsCustom(err) && code >= http.StatusInternalServerError {
		h.newResponse(w, code, response{
			Msg:   msg,
			Error: "unexpected error",
		})
		return
	}

	h.newResponse(w, code, response{
		Msg:   msg,
		Error: err.Error(),
	})
}

func (h *Handler) newResponse(w http.ResponseWriter, code int, data any) {
	h.log.Info("sending response", "status", code, "data", data)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func errIsCustom(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	if errors.Is(err, services.ErrPersonNotExist) {
		return true
	}

	return false
}
