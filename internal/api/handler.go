package api

import (
	"net/http"

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

}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {

}
