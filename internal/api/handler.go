package api

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucent1/rune/internal/store"
)

type Handler struct {
	rune *store.Rune
}

func NewHandler(rune *store.Rune) *Handler {
	return &Handler{
		rune: rune,
	}
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	val, err := io.ReadAll(r.Body)
	if err != nil || val == nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	h.rune.Set(key, val)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	val := h.rune.Get(key)
	w.Write(val)
}
