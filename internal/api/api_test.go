package api

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucent1/rune/internal/config"
	"github.com/lucent1/rune/internal/metrics"
	"github.com/lucent1/rune/internal/store"
	"github.com/prometheus/client_golang/prometheus"
)

func TestHTTPSetAndGet(t *testing.T) {
	rune := store.NewRune()
	cfg := config.Config{MaxKeySize: 256, MaxValueSize: 1048576}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	m := metrics.New(prometheus.NewRegistry())
	router := NewRouter(rune, cfg, logger, m)

	server := httptest.NewServer(router)
	defer server.Close()

	// Set a key
	resp, err := http.PostForm(server.URL+"/set", map[string][]string{
		"key":   {"testkey"},
		"value": {"testvalue"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	// Get it back
	resp, err = http.Get(server.URL + "/get?key=testkey")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "testvalue" {
		t.Errorf("expected testvalue, got %s", body)
	}
}

func TestHTTPGetNotFound(t *testing.T) {
	rune := store.NewRune()
	cfg := config.Config{MaxKeySize: 256, MaxValueSize: 1048576}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	m := metrics.New(prometheus.NewRegistry())
	router := NewRouter(rune, cfg, logger, m)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/get?key=nonexistent")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
