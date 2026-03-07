package aladhan

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

func TestClientGetByCoordinatesParsesSchedule(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "aladhan", "timings_istanbul.json"))
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(&http.Client{Timeout: time.Second}, server.URL)
	schedule, err := client.GetByCoordinates(context.Background(), 41.01384, 28.94966, time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetByCoordinates returned error: %v", err)
	}

	if schedule.MethodID != 13 {
		t.Fatalf("MethodID = %d, want 13", schedule.MethodID)
	}
	if !schedule.RamadanActive {
		t.Fatal("RamadanActive = false, want true")
	}
	if got := schedule.MaghribAt.Format(time.RFC3339); got != "2026-03-08T19:10:00+03:00" {
		t.Fatalf("MaghribAt = %q, want 2026-03-08T19:10:00+03:00", got)
	}
}

func TestClientGetByCoordinatesRejectsMalformedTime(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "aladhan", "timings_bad_time.json"))
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(&http.Client{Timeout: time.Second}, server.URL)
	_, err := client.GetByCoordinates(context.Background(), 41.01384, 28.94966, time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC))
	if err == nil {
		t.Fatal("expected malformed time error")
	}

	cliErr := app.AsCLIError(err)
	if cliErr.ExitCode != app.ExitFailure {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, app.ExitFailure)
	}
}
