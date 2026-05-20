package main

import (
	"log"
	"net/http"

	"github.com/lucent1/rune/internal/api"
	"github.com/lucent1/rune/internal/store"
)

func main() {
	rune := store.NewRune()

	router := api.NewRouter(rune)

	log.Printf("Server running on: 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
