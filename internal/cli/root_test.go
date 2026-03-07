package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

func TestCLIGoldenOutputs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		args       []string
		wantFile   string
		wantExit   int
		wantStream string
	}{
		{
			name:       "help",
			args:       []string{"times", "get", "--help"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_help.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "json",
			args:       []string{"times", "get", "--query", "Istanbul", "--json"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_json.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "human",
			args:       []string{"times", "get", "--query", "Istanbul"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_human.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "quiet",
			args:       []string{"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--quiet"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_countdown_quiet.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, Dependencies{
				Resolver: cliFakeResolver{},
				Provider: cliFakeProvider{},
				Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
			}, tc.args...)

			if exitCode != tc.wantExit {
				t.Fatalf("exitCode = %d, want %d", exitCode, tc.wantExit)
			}

			got := stdout
			if tc.wantStream == "stderr" {
				got = stderr
			}

			want, err := os.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", tc.wantFile, err)
			}

			normalizedGot := strings.TrimRight(got, "\n")
			normalizedWant := strings.TrimRight(string(want), "\n")
			if normalizedGot != normalizedWant {
				t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", tc.name, got, string(want))
			}
		})
	}
}

func TestCLIUsageErrorsHonorJSONOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
	}{
		{
			name: "runtime validation error",
			args: []string{"times", "get", "--json"},
		},
		{
			name: "flag parse error",
			args: []string{"times", "get", "--json", "--badflag"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, Dependencies{
				Resolver: cliFakeResolver{},
				Provider: cliFakeProvider{},
				Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
			}, tc.args...)

			if exitCode != app.ExitUsage {
				t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitUsage)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var payload map[string]any
			if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
				t.Fatalf("Unmarshal(stdout): %v\nstdout=%q", err, stdout)
			}

			want := map[string]any{
				"ok":             false,
				"exit_code":      float64(app.ExitUsage),
				"error_type":     "usage_error",
				"input_received": "",
			}
			for key, wantValue := range want {
				if !reflect.DeepEqual(payload[key], wantValue) {
					t.Fatalf("payload[%q] = %#v, want %#v", key, payload[key], wantValue)
				}
			}
			if payload["message"] == "" {
				t.Fatal("payload[\"message\"] is empty")
			}
			if payload["suggestion"] == "" {
				t.Fatal("payload[\"suggestion\"] is empty")
			}
		})
	}
}

func TestCLIRejectsUnsupportedFieldBeforeProviderCall(t *testing.T) {
	t.Parallel()

	provider := &cliSpyProvider{}
	stdout, stderr, exitCode := executeTestCommand(t, Dependencies{
		Resolver: cliFakeResolver{},
		Provider: provider,
		Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
	}, "times", "get", "--lat", "41.01384", "--lon", "28.94966", "--field", "bogus")

	if exitCode != app.ExitUsage {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitUsage)
	}
	if provider.calls != 0 {
		t.Fatalf("provider.calls = %d, want 0", provider.calls)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `unsupported field "bogus"`) {
		t.Fatalf("stderr = %q, want unsupported field error", stderr)
	}
}

func TestCLIAcceptsRamadanActiveField(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(t, Dependencies{
		Resolver: cliFakeResolver{},
		Provider: cliFakeProvider{},
		Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
	}, "times", "get", "--query", "Istanbul", "--field", "ramadan_active", "--quiet")

	if exitCode != app.ExitSuccess {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if stdout != "true\n" {
		t.Fatalf("stdout = %q, want %q", stdout, "true\n")
	}
}

func executeTestCommand(t *testing.T, deps Dependencies, args ...string) (string, string, int) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	deps.Stdout = &stdout
	deps.Stderr = &stderr

	cmd := NewRootCmd(deps)
	cmd.SetArgs(args)

	exitCode := app.ExitSuccess
	if err := cmd.Execute(); err != nil {
		exitCode = renderCommandError(&stdout, &stderr, isJSONEnabled(cmd), err)
	}

	return stdout.String(), stderr.String(), exitCode
}

type cliFakeResolver struct{}

func (cliFakeResolver) Search(ctx context.Context, query, countryCode string, limit int) ([]app.Location, error) {
	return []app.Location{{
		Name:        "Istanbul",
		Country:     "Türkiye",
		CountryCode: "TR",
		Admin1:      "Istanbul",
		Latitude:    41.01384,
		Longitude:   28.94966,
		Timezone:    "Europe/Istanbul",
	}}, nil
}

type cliFakeProvider struct{}

func (cliFakeProvider) GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (app.DaySchedule, error) {
	tz, _ := time.LoadLocation("Europe/Istanbul")
	return app.DaySchedule{
		Latitude:      latitude,
		Longitude:     longitude,
		Timezone:      "Europe/Istanbul",
		Date:          time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, tz),
		ImsakAt:       time.Date(date.Year(), date.Month(), date.Day(), 5, 46, 0, 0, tz),
		FajrAt:        time.Date(date.Year(), date.Month(), date.Day(), 5, 56, 0, 0, tz),
		SunriseAt:     time.Date(date.Year(), date.Month(), date.Day(), 7, 21, 0, 0, tz),
		DhuhrAt:       time.Date(date.Year(), date.Month(), date.Day(), 13, 20, 0, 0, tz),
		AsrAt:         time.Date(date.Year(), date.Month(), date.Day(), 16, 31, 0, 0, tz),
		MaghribAt:     time.Date(date.Year(), date.Month(), date.Day(), 19, 9, 0, 0, tz),
		SunsetAt:      time.Date(date.Year(), date.Month(), date.Day(), 19, 9, 0, 0, tz),
		IshaAt:        time.Date(date.Year(), date.Month(), date.Day(), 20, 28, 0, 0, tz),
		MethodID:      13,
		MethodName:    "Diyanet İşleri Başkanlığı, Turkey (experimental)",
		Source:        "aladhan:method=13",
		RamadanActive: true,
	}, nil
}

type cliSpyProvider struct {
	calls int
}

func (p *cliSpyProvider) GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (app.DaySchedule, error) {
	p.calls++
	return app.DaySchedule{}, nil
}

type cliFixedClock struct {
	now time.Time
}

func (c cliFixedClock) Now() time.Time {
	return c.now
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()

	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("LoadLocation(%q): %v", name, err)
	}

	return location
}
