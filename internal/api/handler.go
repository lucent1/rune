package api

import (
	"net/http"

	"github.com/lucent1/rune/internal/store"
)

const (
	maxKeySize   = 256
	maxValueSize = 1_000_000
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
	r.Body = http.MaxBytesReader(w, r.Body, maxValueSize+maxKeySize+100)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form or too large", http.StatusBadRequest)
		return
	}

	key := r.PostForm.Get("key")
	val := r.PostForm.Get("value")

	if len(key) == 0 {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if len(key) > maxKeySize {
		http.Error(w, "key too large", http.StatusBadRequest)
		return
	}

	if len(val) == 0 {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	if len(val) > maxValueSize {
		http.Error(w, "value too large", http.StatusBadRequest)
		return
	}

	h.rune.Set(key, []byte(val))

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if len(key) > 256 {
		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	val := h.rune.Get(key)

	w.Header().Set("Content-Type", "text/plain")
	w.Write(val)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if len(key) > 256 {
		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	h.rune.Delete(key)
	w.WriteHeader(http.StatusAccepted)
}
