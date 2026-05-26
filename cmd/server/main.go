package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/lucent1/rune/internal/api"
	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/store"
)

func main() {
	cfg := config.LoadConfig()

	logger := initLogger(cfg.LogLevel)

	rune := store.NewRune()

	router := api.NewRouter(rune, cfg, logger)

	addr := fmt.Sprintf(":%d", cfg.Port)

	log.Printf("Server running on %s", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}

func initLogger(level string) *slog.Logger {
	var lvl slog.Level

	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "error":
		lvl = slog.LevelError
	case "warn":
		lvl = slog.LevelWarn
	default:
		lvl = slog.LevelInfo
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		}),
	)

	return logger
}
