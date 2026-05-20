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
	const (
		maxKeySize   = 256
		maxValueSize = 1_000_000
	)

	key := chi.URLParam(r, "key")

	if len(key) == 0 {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if len(key) > maxKeySize {
		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxValueSize)

	val, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "value too large or invalid body", http.StatusRequestEntityTooLarge)
		return
	}

	h.rune.Set(key, val)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if len(key) > 256 {
		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	val := h.rune.Get(key)
	w.Write(val)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if len(key) > 256 {
		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	h.rune.Delete(key)
	w.WriteHeader(http.StatusAccepted)
}
