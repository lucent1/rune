package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucent1/rune/internal/store"
)

func NewRouter(rune *store.Rune) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(rune)

	r.Put("/key/{key}", h.Set)
	r.Get("/key/{key}", h.Get)
	r.Delete("/key/{key}", h.Delete)

	return r
}
