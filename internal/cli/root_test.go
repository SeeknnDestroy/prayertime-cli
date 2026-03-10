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
			name:       "root_help",
			args:       []string{"--help"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "root_help.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_get_help",
			args:       []string{"times", "get", "--help"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_help.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_countdown_help",
			args:       []string{"times", "countdown", "--help"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_countdown_help.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_get_json_detailed",
			args:       []string{"times", "get", "--query", "Istanbul", "--output", "json"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_json_detailed.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_get_json_concise",
			args:       []string{"times", "get", "--query", "Istanbul", "--output", "json", "--view", "concise"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_json_concise.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_countdown_json_concise",
			args:       []string{"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--output", "json"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_countdown_json_concise.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_countdown_json_detailed",
			args:       []string{"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--output", "json", "--view", "detailed"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_countdown_json_detailed.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_get_human",
			args:       []string{"times", "get", "--query", "Istanbul"},
			wantFile:   filepath.Join("..", "..", "testdata", "golden", "times_get_human.txt"),
			wantExit:   app.ExitSuccess,
			wantStream: "stdout",
		},
		{
			name:       "times_countdown_quiet",
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

			stdout, stderr, exitCode := executeTestCommand(t, testDependencies(t), tc.args...)
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
			args: []string{"times", "get", "--output", "json"},
		},
		{
			name: "flag parse error",
			args: []string{"times", "get", "--output", "json", "--badflag"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, testDependencies(t), tc.args...)
			if exitCode != app.ExitUsage {
				t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitUsage)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			payload := decodeJSON(t, stdout)
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

func TestCLIAcceptsRamadanActiveFieldWithOutputValue(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "get", "--query", "Istanbul", "--field", "ramadan_active", "--output", "value",
	)

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

func TestCLIQuietAliasMatchesOutputValue(t *testing.T) {
	t.Parallel()

	quietStdout, quietStderr, quietExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--quiet",
	)
	valueStdout, valueStderr, valueExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--field", "seconds_remaining", "--output", "value",
	)

	if quietExitCode != app.ExitSuccess || valueExitCode != app.ExitSuccess {
		t.Fatalf("quiet/value exit codes = %d/%d, want success", quietExitCode, valueExitCode)
	}
	if quietStderr != "" || valueStderr != "" {
		t.Fatalf("stderr should be empty, got quiet=%q value=%q", quietStderr, valueStderr)
	}
	if quietStdout != valueStdout {
		t.Fatalf("quiet stdout = %q, value stdout = %q", quietStdout, valueStdout)
	}
}

func TestCLIJSONAliasMatchesOutputJSON(t *testing.T) {
	t.Parallel()

	jsonAliasStdout, jsonAliasStderr, jsonAliasExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "get", "--query", "Istanbul", "--json",
	)
	outputStdout, outputStderr, outputExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "get", "--query", "Istanbul", "--output", "json",
	)

	if jsonAliasExitCode != app.ExitSuccess || outputExitCode != app.ExitSuccess {
		t.Fatalf("json alias/output exit codes = %d/%d, want success", jsonAliasExitCode, outputExitCode)
	}
	if jsonAliasStderr != "" || outputStderr != "" {
		t.Fatalf("stderr should be empty, got alias=%q output=%q", jsonAliasStderr, outputStderr)
	}
	if jsonAliasStdout != outputStdout {
		t.Fatalf("json alias stdout = %q, output stdout = %q", jsonAliasStdout, outputStdout)
	}
}

func TestCLIHelpShowsAgentShortcuts(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "root help shows json",
			args: []string{"--help"},
			contains: []string{
				"--json",
				"Shortcut for --output json",
			},
		},
		{
			name: "times get help shows quiet",
			args: []string{"times", "get", "--help"},
			contains: []string{
				"--quiet",
				"Shortcut for --output value",
			},
		},
		{
			name: "times countdown help shows quiet",
			args: []string{"times", "countdown", "--help"},
			contains: []string{
				"--quiet",
				"Shortcut for --output value",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, testDependencies(t), tc.args...)
			if exitCode != app.ExitSuccess {
				t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitSuccess)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			for _, needle := range tc.contains {
				if !strings.Contains(stdout, needle) {
					t.Fatalf("stdout = %q, want %q", stdout, needle)
				}
			}
		})
	}
}

func TestCLIJSONFalsePreservesTextOutput(t *testing.T) {
	t.Parallel()

	defaultStdout, defaultStderr, defaultExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"version",
	)
	falseStdout, falseStderr, falseExitCode := executeTestCommand(
		t,
		testDependencies(t),
		"version", "--json=false",
	)

	if defaultExitCode != app.ExitSuccess || falseExitCode != app.ExitSuccess {
		t.Fatalf("default/false exit codes = %d/%d, want success", defaultExitCode, falseExitCode)
	}
	if defaultStderr != "" || falseStderr != "" {
		t.Fatalf("stderr should be empty, got default=%q false=%q", defaultStderr, falseStderr)
	}
	if falseStdout != defaultStdout {
		t.Fatalf("json=false stdout = %q, default stdout = %q", falseStdout, defaultStdout)
	}
}

func TestCLIJSONFalseErrorsRenderToStderr(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "get", "--json=false", "--badflag",
	)

	if exitCode != app.ExitUsage {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitUsage)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "unknown flag: --badflag") {
		t.Fatalf("stderr = %q, want flag parse error", stderr)
	}
}

func TestCLIRejectsExplicitEmptyOutputMode(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
	}{
		{
			name: "version",
			args: []string{"version", "--output", ""},
		},
		{
			name: "times get",
			args: []string{"times", "get", "--output", ""},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, testDependencies(t), tc.args...)
			if exitCode != app.ExitUsage {
				t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitUsage)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, `unsupported output mode ""`) {
				t.Fatalf("stderr = %q, want unsupported output mode error", stderr)
			}
		})
	}
}

func TestCLIJSONErrorsIncludeDetails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		deps    Dependencies
		args    []string
		wantKey string
	}{
		{
			name:    "invalid field",
			deps:    testDependencies(t),
			args:    []string{"times", "get", "--lat", "41.01384", "--lon", "28.94966", "--field", "bogus", "--output", "json"},
			wantKey: "valid_fields",
		},
		{
			name:    "invalid target",
			deps:    testDependencies(t),
			args:    []string{"times", "countdown", "--query", "Istanbul", "--target", "bogus", "--output", "json"},
			wantKey: "valid_targets",
		},
		{
			name:    "missing location input",
			deps:    testDependencies(t),
			args:    []string{"times", "get", "--output", "json"},
			wantKey: "required_one_of",
		},
		{
			name: "ambiguous location",
			deps: Dependencies{
				Resolver: cliAmbiguousResolver{},
				Provider: cliFakeProvider{},
				Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
			},
			args:    []string{"times", "get", "--query", "Springfield", "--output", "json"},
			wantKey: "candidates",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, exitCode := executeTestCommand(t, tc.deps, tc.args...)
			if exitCode == app.ExitSuccess {
				t.Fatal("expected non-success exit code")
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			payload := decodeJSON(t, stdout)
			details, ok := payload["details"].(map[string]any)
			if !ok {
				t.Fatalf("details = %#v, want object", payload["details"])
			}
			if _, ok := details[tc.wantKey]; !ok {
				t.Fatalf("details[%q] missing in %#v", tc.wantKey, details)
			}
		})
	}
}

func TestCLIAmbiguousLocationRendersCandidatesToStderr(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(t, Dependencies{
		Resolver: cliAmbiguousResolver{},
		Provider: cliFakeProvider{},
		Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
	}, "times", "get", "--query", "Springfield")

	if exitCode != app.ExitNotFound {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitNotFound)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "Candidates:\n1. Springfield, Illinois, United States") {
		t.Fatalf("stderr = %q, want numbered candidate list", stderr)
	}
	if strings.Count(stderr, "Springfield, Missouri, United States") != 1 {
		t.Fatalf("stderr = %q, want deduped candidate list", stderr)
	}
}

func TestCLILocationsSearchJSONIncludesDisplayName(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(
		t,
		testDependencies(t),
		"locations", "search", "--query", "Istanbul", "--output", "json",
	)

	if exitCode != app.ExitSuccess {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	payload := decodeJSON(t, stdout)
	results, ok := payload["results"].([]any)
	if !ok || len(results) == 0 {
		t.Fatalf("results = %#v, want non-empty array", payload["results"])
	}
	first, ok := results[0].(map[string]any)
	if !ok {
		t.Fatalf("results[0] = %#v, want object", results[0])
	}
	if first["display_name"] != "Istanbul, Türkiye" {
		t.Fatalf("display_name = %#v, want %q", first["display_name"], "Istanbul, Türkiye")
	}
}

func TestCLICountdownFieldJSONReturnsSingleKeyObject(t *testing.T) {
	t.Parallel()

	stdout, stderr, exitCode := executeTestCommand(
		t,
		testDependencies(t),
		"times", "countdown", "--query", "Istanbul", "--target", "iftar", "--field", "seconds_remaining", "--output", "json",
	)

	if exitCode != app.ExitSuccess {
		t.Fatalf("exitCode = %d, want %d", exitCode, app.ExitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	payload := decodeJSON(t, stdout)
	if len(payload) != 1 {
		t.Fatalf("len(payload) = %d, want 1", len(payload))
	}
	if _, ok := payload["seconds_remaining"]; !ok {
		t.Fatalf("payload = %#v, want seconds_remaining only", payload)
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

func decodeJSON(t *testing.T, value string) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		t.Fatalf("Unmarshal(%q): %v", value, err)
	}
	return payload
}

func testDependencies(t *testing.T) Dependencies {
	t.Helper()

	return Dependencies{
		Resolver: cliFakeResolver{},
		Provider: cliFakeProvider{},
		Clock:    cliFixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, mustLocation(t, "Europe/Istanbul"))},
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

type cliAmbiguousResolver struct{}

func (cliAmbiguousResolver) Search(ctx context.Context, query, countryCode string, limit int) ([]app.Location, error) {
	return []app.Location{
		{Name: "Springfield", Country: "United States", CountryCode: "US", Admin1: "Illinois", Latitude: 39.78, Longitude: -89.64, Timezone: "America/Chicago"},
		{Name: "Springfield", Country: "United States", CountryCode: "US", Admin1: "Missouri", Latitude: 37.20, Longitude: -93.29, Timezone: "America/Chicago"},
		{Name: "Springfield", Country: "United States", CountryCode: "US", Admin1: "Missouri", Latitude: 37.20, Longitude: -93.29, Timezone: "America/Chicago"},
	}, nil
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
