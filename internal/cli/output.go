package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

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
		_ = writeJSON(stdout, payload)
		return cliErr.ExitCode
	}

	_, _ = fmt.Fprintln(stderr, cliErr.Message)
	if cliErr.Suggestion != "" {
		_, _ = fmt.Fprintln(stderr, cliErr.Suggestion)
	}
	return cliErr.ExitCode
}

func writeJSON(out io.Writer, payload any) error {
	encoder := json.NewEncoder(out)
	return encoder.Encode(payload)
}

func writeTimesHuman(out io.Writer, response app.TimesResponse, field string, quiet bool) error {
	if field != "" {
		value, resolvedField, err := selectTimesField(response, field)
		if err != nil {
			return err
		}
		if quiet {
			_, err = fmt.Fprintln(out, value)
			return err
		}
		_, err = fmt.Fprintf(out, "%s: %s\n", resolvedField, value)
		return err
	}

	_, err := fmt.Fprintf(
		out,
		"Location: %s\nTimezone: %s\nDate: %s\nImsak: %s\nFajr: %s\nSunrise: %s\nDhuhr: %s\nAsr: %s\nMaghrib: %s\nSunset: %s\nIsha: %s\nMethod: %d %s\nSource: %s\nRamadan active: %t\n",
		response.LocationName,
		response.Timezone,
		response.Date,
		response.ImsakAt,
		response.FajrAt,
		response.SunriseAt,
		response.DhuhrAt,
		response.AsrAt,
		response.MaghribAt,
		response.SunsetAt,
		response.IshaAt,
		response.MethodID,
		response.MethodName,
		response.Source,
		response.RamadanActive,
	)
	return err
}

func writeCountdownHuman(out io.Writer, response app.CountdownResponse, quiet bool) error {
	if quiet {
		_, err := fmt.Fprintln(out, response.SecondsRemaining)
		return err
	}

	_, err := fmt.Fprintf(
		out,
		"Location: %s\nTimezone: %s\nDate: %s\nTarget: %s\nTarget at: %s\nSeconds remaining: %d\nMinutes remaining: %d\nRamadan active: %t\n",
		response.LocationName,
		response.Timezone,
		response.Date,
		response.Target,
		response.TargetAt,
		response.SecondsRemaining,
		response.MinutesRemaining,
		response.RamadanActive,
	)
	return err
}

func resolveTimesField(field string) (string, error) {
	resolved, ok := app.NormalizeField(field)
	if !ok {
		return "", app.NewUsageError(
			fmt.Sprintf("unsupported field %q", field),
			field,
			"Use fields like imsak, fajr, iftar, timezone, method_name, or source.",
		)
	}

	return resolved, nil
}

func selectTimesField(response app.TimesResponse, field string) (string, string, error) {
	resolved, err := resolveTimesField(field)
	if err != nil {
		return "", "", err
	}

	values := map[string]string{
		"location_name":  response.LocationName,
		"latitude":       strconv.FormatFloat(response.Latitude, 'f', 6, 64),
		"longitude":      strconv.FormatFloat(response.Longitude, 'f', 6, 64),
		"timezone":       response.Timezone,
		"date":           response.Date,
		"imsak_at":       response.ImsakAt,
		"fajr_at":        response.FajrAt,
		"sunrise_at":     response.SunriseAt,
		"dhuhr_at":       response.DhuhrAt,
		"asr_at":         response.AsrAt,
		"maghrib_at":     response.MaghribAt,
		"sunset_at":      response.SunsetAt,
		"isha_at":        response.IshaAt,
		"method_id":      strconv.Itoa(response.MethodID),
		"method_name":    response.MethodName,
		"source":         response.Source,
		"ramadan_active": strconv.FormatBool(response.RamadanActive),
	}

	return values[resolved], resolved, nil
}
