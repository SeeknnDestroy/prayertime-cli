package openmeteo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

func TestClientSearchParsesResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "openmeteo", "search_istanbul.json"))
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(&http.Client{Timeout: time.Second}, server.URL)
	results, err := client.Search(context.Background(), "Istanbul", "TR", 5)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].Name != "Istanbul" {
		t.Fatalf("results[0].Name = %q, want Istanbul", results[0].Name)
	}
}

func TestClientSearchReturnsDecodeErrorForMalformedPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "openmeteo", "search_malformed.json"))
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(&http.Client{Timeout: time.Second}, server.URL)
	_, err := client.Search(context.Background(), "Istanbul", "TR", 5)
	if err == nil {
		t.Fatal("expected decode error")
	}

	cliErr := app.AsCLIError(err)
	if cliErr.ExitCode != app.ExitFailure {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, app.ExitFailure)
	}
}
