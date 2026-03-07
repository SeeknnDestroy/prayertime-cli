package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
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

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd := NewRootCmd(Dependencies{
				Stdout:   &stdout,
				Stderr:   &stderr,
				Resolver: cliFakeResolver{},
				Provider: cliFakeProvider{},
				Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
			})
			cmd.SetArgs(tc.args)

			exitCode := app.ExitSuccess
			if err := cmd.Execute(); err != nil {
				exitCode = renderCommandError(&stdout, &stderr, isJSONEnabled(cmd), err)
			}

			if exitCode != tc.wantExit {
				t.Fatalf("exitCode = %d, want %d", exitCode, tc.wantExit)
			}

			got := stdout.String()
			if tc.wantStream == "stderr" {
				got = stderr.String()
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
