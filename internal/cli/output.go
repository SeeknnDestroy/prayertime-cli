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

func timesFieldEntry(response app.TimesResponse, field string) (fieldEntry, error) {
	resolved, ok := app.NormalizeTimesField(field)
	if !ok {
		return fieldEntry{}, app.NewUsageError(
			fmt.Sprintf("unsupported field %q", field),
			field,
			"Use a canonical prayer-times field such as maghrib_at, timezone, method_name, or source.",
		).WithDetails(app.ErrorDetails{ValidFields: app.ValidTimesFields()})
	}

	for _, entry := range timesEntries(response, viewDetailed) {
		if entry.Key == resolved {
			return entry, nil
		}
	}

	return fieldEntry{}, app.NewInternalError("failed to resolve field from response", resolved, "", nil)
}

func countdownFieldEntry(response app.CountdownResponse, field string) (fieldEntry, error) {
	resolved, ok := app.NormalizeCountdownField(field)
	if !ok {
		return fieldEntry{}, app.NewUsageError(
			fmt.Sprintf("unsupported field %q", field),
			field,
			"Use a canonical countdown field such as seconds_remaining, target_at, maghrib_at, or method_name.",
		).WithDetails(app.ErrorDetails{ValidFields: app.ValidCountdownFields()})
	}

	for _, entry := range countdownEntries(response, viewDetailed) {
		if entry.Key == resolved {
			return entry, nil
		}
	}

	return fieldEntry{}, app.NewInternalError("failed to resolve field from response", resolved, "", nil)
}

func validateTimesField(field string) error {
	if field == "" {
		return nil
	}

	if _, ok := app.NormalizeTimesField(field); ok {
		return nil
	}

	return app.NewUsageError(
		fmt.Sprintf("unsupported field %q", field),
		field,
		"Use a canonical prayer-times field such as maghrib_at, timezone, method_name, or source.",
	).WithDetails(app.ErrorDetails{ValidFields: app.ValidTimesFields()})
}

func validateCountdownField(field string) error {
	if field == "" {
		return nil
	}

	if _, ok := app.NormalizeCountdownField(field); ok {
		return nil
	}

	return app.NewUsageError(
		fmt.Sprintf("unsupported field %q", field),
		field,
		"Use a canonical countdown field such as seconds_remaining, target_at, maghrib_at, or method_name.",
	).WithDetails(app.ErrorDetails{ValidFields: app.ValidCountdownFields()})
}

func timesEntries(response app.TimesResponse, view viewMode) []fieldEntry {
	entries := []fieldEntry{
		{Key: "location_name", Label: "Location", Value: response.LocationName},
		{Key: "timezone", Label: "Timezone", Value: response.Timezone},
		{Key: "date", Label: "Date", Value: response.Date},
		{Key: "imsak_at", Label: "Imsak", Value: response.ImsakAt},
		{Key: "fajr_at", Label: "Fajr", Value: response.FajrAt},
		{Key: "sunrise_at", Label: "Sunrise", Value: response.SunriseAt},
		{Key: "dhuhr_at", Label: "Dhuhr", Value: response.DhuhrAt},
		{Key: "asr_at", Label: "Asr", Value: response.AsrAt},
		{Key: "maghrib_at", Label: "Maghrib", Value: response.MaghribAt},
		{Key: "sunset_at", Label: "Sunset", Value: response.SunsetAt},
		{Key: "isha_at", Label: "Isha", Value: response.IshaAt},
		{Key: "ramadan_active", Label: "Ramadan active", Value: response.RamadanActive},
	}

	if view == viewDetailed {
		entries = append(entries,
			fieldEntry{Key: "latitude", Label: "Latitude", Value: response.Latitude},
			fieldEntry{Key: "longitude", Label: "Longitude", Value: response.Longitude},
			fieldEntry{Key: "method_id", Label: "Method ID", Value: response.MethodID},
			fieldEntry{Key: "method_name", Label: "Method", Value: response.MethodName},
			fieldEntry{Key: "source", Label: "Source", Value: response.Source},
		)
	}

	return entries
}

func countdownEntries(response app.CountdownResponse, view viewMode) []fieldEntry {
	entries := []fieldEntry{
		{Key: "location_name", Label: "Location", Value: response.LocationName},
		{Key: "timezone", Label: "Timezone", Value: response.Timezone},
		{Key: "date", Label: "Date", Value: response.Date},
		{Key: "target", Label: "Target", Value: response.Target},
		{Key: "target_at", Label: "Target at", Value: response.TargetAt},
		{Key: "seconds_remaining", Label: "Seconds remaining", Value: response.SecondsRemaining},
		{Key: "minutes_remaining", Label: "Minutes remaining", Value: response.MinutesRemaining},
	}

	if view == viewDetailed {
		entries = append(entries,
			fieldEntry{Key: "imsak_at", Label: "Imsak", Value: response.ImsakAt},
			fieldEntry{Key: "fajr_at", Label: "Fajr", Value: response.FajrAt},
			fieldEntry{Key: "sunrise_at", Label: "Sunrise", Value: response.SunriseAt},
			fieldEntry{Key: "dhuhr_at", Label: "Dhuhr", Value: response.DhuhrAt},
			fieldEntry{Key: "asr_at", Label: "Asr", Value: response.AsrAt},
			fieldEntry{Key: "maghrib_at", Label: "Maghrib", Value: response.MaghribAt},
			fieldEntry{Key: "sunset_at", Label: "Sunset", Value: response.SunsetAt},
			fieldEntry{Key: "isha_at", Label: "Isha", Value: response.IshaAt},
			fieldEntry{Key: "ramadan_active", Label: "Ramadan active", Value: response.RamadanActive},
			fieldEntry{Key: "latitude", Label: "Latitude", Value: response.Latitude},
			fieldEntry{Key: "longitude", Label: "Longitude", Value: response.Longitude},
			fieldEntry{Key: "method_id", Label: "Method ID", Value: response.MethodID},
			fieldEntry{Key: "method_name", Label: "Method", Value: response.MethodName},
			fieldEntry{Key: "source", Label: "Source", Value: response.Source},
		)
	}

	return entries
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
