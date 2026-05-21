package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucent1/rune/internal/metrics"
	"github.com/lucent1/rune/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(rune *store.Rune) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	m := metrics.New()
	r.Use(m.MiddleWare)

	h := NewHandler(rune)

	r.Post("/set", h.Set)
	r.Get("/get", h.Get)
	r.Delete("/delete", h.Delete)

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	return r
}
