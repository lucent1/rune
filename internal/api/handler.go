package api

import (
	"net/http"

	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/store"
)

type Handler struct {
	rune *store.Rune
	cfg  config.Config
}

func NewHandler(rune *store.Rune, cfg config.Config) *Handler {
	return &Handler{
		rune: rune,
		cfg:  cfg,
	}
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(h.cfg.MaxValueSize+h.cfg.MaxKeySize+100))

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

	if len(key) > h.cfg.MaxKeySize {
		http.Error(w, "key too large", http.StatusBadRequest)
		return
	}

	if len(val) == 0 {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	if len(val) > h.cfg.MaxValueSize {
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
