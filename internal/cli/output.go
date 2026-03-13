package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

type fieldEntry struct {
	Key   string
	Label string
	Value any
}

func renderCommandError(stdout io.Writer, stderr io.Writer, useJSON bool, err error) int {
	cliErr := app.AsCLIError(err)
	if useJSON {
		payload := map[string]any{
			"ok":             false,
			"exit_code":      cliErr.ExitCode,
			"error_type":     cliErr.ErrorType,
			"message":        cliErr.Message,
			"input_received": cliErr.InputReceived,
			"suggestion":     cliErr.Suggestion,
		}
		if cliErr.Details != nil {
			payload["details"] = cliErr.Details
		}
		_ = writeJSON(stdout, payload)
		return cliErr.ExitCode
	}

	_, _ = fmt.Fprintln(stderr, cliErr.Message)
	if cliErr.Details != nil && len(cliErr.Details.Candidates) > 0 {
		_, _ = fmt.Fprintln(stderr, "Candidates:")
		for index, candidate := range cliErr.Details.Candidates {
			_, _ = fmt.Fprintf(
				stderr,
				"%d. %s [%0.6f, %0.6f] %s\n",
				index+1,
				candidate.DisplayName,
				candidate.Latitude,
				candidate.Longitude,
				candidate.Timezone,
			)
		}
	}
	if cliErr.Suggestion != "" {
		_, _ = fmt.Fprintln(stderr, cliErr.Suggestion)
	}
	return cliErr.ExitCode
}

func writeJSON(out io.Writer, payload any) error {
	encoder := json.NewEncoder(out)
	return encoder.Encode(payload)
}

func writeTimesOutput(out io.Writer, response app.TimesResponse, mode outputMode, view viewMode, field string) error {
	if field != "" {
		entry, err := timesFieldEntry(response, field)
		return writeFieldOutput(out, mode, entry, err)
	}

	switch mode {
	case outputJSON:
		return writeJSON(out, entriesToMap(timesEntries(response, view)))
	case outputText:
		return writeTextEntries(out, timesEntries(response, view))
	default:
		return app.NewUsageError(
			fmt.Sprintf("unsupported output mode %q for times get", mode),
			string(mode),
			"Use --output text, --output json, or --output value with --field.",
		)
	}
}

func writeCountdownOutput(out io.Writer, response app.CountdownResponse, mode outputMode, view viewMode, field string) error {
	if field != "" {
		entry, err := countdownFieldEntry(response, field)
		return writeFieldOutput(out, mode, entry, err)
	}

	switch mode {
	case outputJSON:
		return writeJSON(out, entriesToMap(countdownEntries(response, view)))
	case outputText:
		return writeTextEntries(out, countdownEntries(response, view))
	default:
		return app.NewUsageError(
			fmt.Sprintf("unsupported output mode %q for times countdown", mode),
			string(mode),
			"Use --output text, --output json, or --output value with --field.",
		)
	}
}

func writeFieldOutput(out io.Writer, mode outputMode, entry fieldEntry, err error) error {
	if err != nil {
		return err
	}

	switch mode {
	case outputText:
		_, err = fmt.Fprintf(out, "%s: %s\n", entry.Key, stringifyValue(entry.Value))
		return err
	case outputJSON:
		return writeJSON(out, map[string]any{entry.Key: entry.Value})
	case outputValue:
		_, err = fmt.Fprintln(out, stringifyValue(entry.Value))
		return err
	default:
		return app.NewUsageError(
			fmt.Sprintf("unsupported output mode %q", mode),
			string(mode),
			"Use --output text, --output json, or --output value.",
		)
	}
}

func writeTextEntries(out io.Writer, entries []fieldEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintf(out, "%s: %s\n", entry.Label, stringifyValue(entry.Value)); err != nil {
			return err
		}
	}
	return nil
}

func entriesToMap(entries []fieldEntry) map[string]any {
	payload := make(map[string]any, len(entries))
	for _, entry := range entries {
		payload[entry.Key] = entry.Value
	}
	return payload
}

func stringifyValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', 6, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprint(typed)
	}
}
