package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lucent1/rune/internal/api"
	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/store"
)

func main() {
	cfg := config.LoadConfig()

	rune := store.NewRune()

	router := api.NewRouter(rune, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)

	log.Printf("Server running on %s", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
