package httpx

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

func TestGetJSONDoesNotRetryPermanentTransportError(t *testing.T) {
	requestCalls := 0
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requestCalls++
			return nil, errors.New("connection refused")
		}),
	}

	var payload map[string]any
	err := GetJSON(context.Background(), client, "https://example.com", &payload)
	if err == nil {
		t.Fatal("expected transport error")
	}

	cliErr := app.AsCLIError(err)
	if cliErr.ExitCode != app.ExitNetwork {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, app.ExitNetwork)
	}
	if cliErr.ErrorType != "network_error" {
		t.Fatalf("ErrorType = %q, want network_error", cliErr.ErrorType)
	}
	if requestCalls != 1 {
		t.Fatalf("requestCalls = %d, want 1", requestCalls)
	}
	if cliErr.Message != "upstream request failed before a response was received" {
		t.Fatalf("Message = %q, want permanent network failure message", cliErr.Message)
	}
}

func TestGetJSONRetriesTimeoutTransportErrors(t *testing.T) {
	originalWaitBeforeRetry := waitBeforeRetry
	waitCalls := 0
	waitBeforeRetry = func(ctx context.Context, duration time.Duration) error {
		waitCalls++
		return nil
	}
	defer func() {
		waitBeforeRetry = originalWaitBeforeRetry
	}()

	requestCalls := 0
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requestCalls++
			return nil, fakeTimeoutError{}
		}),
	}

	var payload map[string]any
	err := GetJSON(context.Background(), client, "https://example.com", &payload)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	cliErr := app.AsCLIError(err)
	if cliErr.ExitCode != app.ExitNetwork {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, app.ExitNetwork)
	}
	if cliErr.ErrorType != "network_timeout" {
		t.Fatalf("ErrorType = %q, want network_timeout", cliErr.ErrorType)
	}
	if requestCalls != 3 {
		t.Fatalf("requestCalls = %d, want 3", requestCalls)
	}
	if waitCalls != 2 {
		t.Fatalf("waitCalls = %d, want 2", waitCalls)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type fakeTimeoutError struct{}

func (fakeTimeoutError) Error() string {
	return "i/o timeout"
}

func (fakeTimeoutError) Timeout() bool {
	return true
}

func (fakeTimeoutError) Temporary() bool {
	return false
}
