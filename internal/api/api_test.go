package api

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
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

	//Delete
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/delete?key=testkey", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("expected 202, got %d", req.Response.StatusCode)
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

	//Get
	resp, err := http.Get(server.URL + "/get?key=nonexistent")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}

	//Delete
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/delete?key=nonexistent", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}

func TestKeyTooLarge(t *testing.T) {
	rune := store.NewRune()
	cfg := config.Config{MaxKeySize: 256, MaxValueSize: 1048576}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	m := metrics.New(prometheus.NewRegistry())
	router := NewRouter(rune, cfg, logger, m)

	server := httptest.NewServer(router)
	defer server.Close()

	//large key
	longKey := strings.Repeat("a", 257)
	longVal := strings.Repeat("b", 1048576+101)

	// Key too large
	resp, err := http.PostForm(server.URL+"/set", map[string][]string{
		"key":   {longKey},
		"value": {"testvalue"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	resp, err = http.PostForm(server.URL+"/set", map[string][]string{
		"key":   {longKey},
		"value": {longVal},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	//long value
	resp, err = http.PostForm(server.URL+"/set", map[string][]string{
		"key":   {"test"},
		"value": {longVal},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	//get key too large
	resp, err = http.Get(server.URL + "/get?key=" + longKey)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	//Delete key too large
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/delete?key="+longKey, nil)

	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

}
