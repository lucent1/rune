package api

import (
	"log/slog"
	"net/http"

	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/store"
)

type Handler struct {
	rune   *store.Rune
	cfg    config.Config
	logger *slog.Logger
}

func NewHandler(rune *store.Rune, cfg config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		rune:   rune,
		cfg:    cfg,
		logger: logger,
	}
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(h.cfg.MaxValueSize+h.cfg.MaxKeySize+100))

	if err := r.ParseForm(); err != nil {
		h.logger.Warn("Invalid form or Too large")

		http.Error(w, "invalid form or too large", http.StatusBadRequest)
		return
	}

	key := r.PostForm.Get("key")
	val := r.PostForm.Get("value")

	h.logger.Debug("set request received",
		"keySize", len(key),
		"valueSize", len(val),
	)

	if len(key) == 0 {
		h.logger.Warn("Key is empty", "key size", len(key))

		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if len(key) > h.cfg.MaxKeySize {
		h.logger.Warn("Key too large",
			"key", len(key),
			"Max allowed", h.cfg.MaxKeySize,
		)

		http.Error(w, "key too large", http.StatusBadRequest)
		return
	}

	if len(val) == 0 {
		h.logger.Warn("Value is empty", "value", len(val))

		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	if len(val) > h.cfg.MaxValueSize {
		h.logger.Warn("Value is too large",
			"value", len(val),
			"Max allowed", h.cfg.MaxValueSize,
		)

		http.Error(w, "value too large", http.StatusBadRequest)
		return
	}

	h.rune.Set(key, []byte(val))

	h.logger.Info("Key stored successfully", "key", key)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	h.logger.Debug("get request received", "key", key)

	if key == "" {
		h.logger.Warn("Key is empty", "keysize", len(key))

		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if len(key) > 256 {
		h.logger.Warn("Key size too large", "keysize", len(key))

		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	val := h.rune.Get(key)

	h.logger.Info("key retrieved", "key", key, "valuesize", len(val))

	w.Header().Set("Content-Type", "text/plain")
	w.Write(val)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if len(key) > 256 {
		h.logger.Warn("Key size too large", "keysize", len(key))

		http.Error(w, "Key size exceeds 256 bytes", http.StatusBadRequest)
		return
	}

	h.rune.Delete(key)
	h.logger.Info("key value deleted")
	w.WriteHeader(http.StatusAccepted)
}
