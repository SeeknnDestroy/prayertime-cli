package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/spf13/cobra"
)

func newTimesCmd(service *app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "times",
		Short: "Fetch daily prayer schedules and countdowns",
	}

	cmd.AddCommand(newTimesGetCmd(service))
	cmd.AddCommand(newTimesCountdownCmd(service))
	return cmd
}

func newTimesGetCmd(service *app.Service) *cobra.Command {
	var query string
	var countryCode string
	var date string
	var field string
	var quiet bool
	var latitude float64
	var longitude float64
	var latitudeSet bool
	var longitudeSet bool

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch prayer times for a location and date",
		Long: strings.TrimSpace(`
Fetch one day's prayer schedule for a resolved place or explicit coordinates.

MVP 1 is stateless, so every call must include --query PLACE or both --lat and --lon. Use --field with --quiet for one scalar value, or --json for the full structured response.
`),
		Example: strings.TrimSpace(`
prayertime-cli times get --query Istanbul
prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
prayertime-cli times get --lat 41.01384 --lon 28.94966 --date 2026-03-07 --json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if isJSONEnabled(cmd) && field != "" {
				return app.NewUsageError("--field cannot be used with --json", field, "Use --quiet with --field for bare value output or drop --json.")
			}
			if isJSONEnabled(cmd) && quiet {
				return app.NewUsageError("--quiet cannot be used with --json", "", "Choose structured JSON or quiet bare output.")
			}
			if quiet && field == "" {
				return app.NewUsageError("--quiet requires --field for 'times get'", "", "Provide --field to choose a single value.")
			}

			resolvedField := ""
			if field != "" {
				var err error
				resolvedField, err = resolveTimesField(field)
				if err != nil {
					return err
				}
			}

			request := app.TimesRequest{
				Query:       query,
				CountryCode: strings.ToUpper(countryCode),
				Date:        date,
			}
			if latitudeSet {
				request.Latitude = &latitude
			}
			if longitudeSet {
				request.Longitude = &longitude
			}

			response, err := service.GetTimes(cmd.Context(), request)
			if err != nil {
				return err
			}

			if isJSONEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), response)
			}

			return writeTimesHuman(cmd.OutOrStdout(), response, resolvedField, quiet)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve. Required unless --lat and --lon are set")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&date, "date", "today", "Date in YYYY-MM-DD format or 'today'")
	cmd.Flags().StringVar(&field, "field", "", "Return one field such as maghrib, iftar, yatsi, timezone, or method_name")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Emit only the selected field value")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate. Use with --lon instead of --query")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate. Use with --lat instead of --query")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		latitudeSet = cmd.Flags().Changed("lat")
		longitudeSet = cmd.Flags().Changed("lon")
	}

	return cmd
}

func newTimesCountdownCmd(service *app.Service) *cobra.Command {
	var query string
	var countryCode string
	var target string
	var at string
	var quiet bool
	var latitude float64
	var longitude float64
	var latitudeSet bool
	var longitudeSet bool

	cmd := &cobra.Command{
		Use:   "countdown",
		Short: "Calculate seconds and minutes remaining until a target prayer",
		Long: strings.TrimSpace(`
Calculate remaining time until the next prayer or a named prayer target.

Use --target next-prayer for generic "next ezan" questions. Canonical targets are English, while Turkish semantic aliases such as iftar, aksam, and yatsi are also accepted.
`),
		Example: strings.TrimSpace(`
prayertime-cli times countdown --query Istanbul --target next-prayer --json
prayertime-cli times countdown --query Istanbul --target iftar --quiet
prayertime-cli times countdown --lat 41.01384 --lon 28.94966 --target asr --at 2026-03-07T12:00:00+03:00 --json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if isJSONEnabled(cmd) && quiet {
				return app.NewUsageError("--quiet cannot be used with --json", "", "Choose structured JSON or quiet bare output.")
			}

			var atValue *time.Time
			if strings.TrimSpace(at) != "" {
				parsed, err := time.Parse(time.RFC3339, at)
				if err != nil {
					return app.NewUsageError(fmt.Sprintf("invalid RFC3339 value %q", at), at, "Use --at 2026-03-07T18:00:00+03:00.")
				}
				atValue = &parsed
			}

			request := app.CountdownRequest{
				Query:       query,
				CountryCode: strings.ToUpper(countryCode),
				Target:      target,
				At:          atValue,
			}
			if latitudeSet {
				request.Latitude = &latitude
			}
			if longitudeSet {
				request.Longitude = &longitude
			}

			response, err := service.GetCountdown(cmd.Context(), request)
			if err != nil {
				return err
			}

			if isJSONEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), response)
			}

			return writeCountdownHuman(cmd.OutOrStdout(), response, quiet)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve. Required unless --lat and --lon are set")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&target, "target", "", "Target prayer. Use next-prayer for generic countdowns; iftar and yatsi are accepted aliases")
	cmd.Flags().StringVar(&at, "at", "", "Optional RFC3339 timestamp to evaluate countdown from")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Emit only remaining seconds")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate. Use with --lon instead of --query")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate. Use with --lat instead of --query")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		latitudeSet = cmd.Flags().Changed("lat")
		longitudeSet = cmd.Flags().Changed("lon")
	}

	return cmd
}
