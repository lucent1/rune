package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucent1/rune/internal/api"
	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/metrics"
	"github.com/lucent1/rune/internal/store"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	cfg := config.LoadConfig()

	logger := initLogger(cfg.LogLevel)

	rune := store.NewRune()

	m := metrics.New(prometheus.NewRegistry())

	router := api.NewRouter(rune, cfg, logger, m)

	addr := fmt.Sprintf(":%d", cfg.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Error("server crashed", "error", err)
		}
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigch)

	sig := <-sigch

	logger.Info("shutdown signal received", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	} else {
		logger.Info("server shutdown complete")
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
