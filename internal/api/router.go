package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/metrics"
	"github.com/lucent1/rune/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(rune *store.Rune, cfg config.Config, logger *slog.Logger, metric *metrics.Metrics) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(metric.MiddleWare)

	h := NewHandler(rune, cfg, logger)

	r.Post("/set", h.Set)
	r.Get("/get", h.Get)
	r.Delete("/delete", h.Delete)

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	return r
}
