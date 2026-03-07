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
		Short: "Fetch daily prayer times and countdowns",
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
		Example: strings.TrimSpace(`
prayertime-cli times get --query Istanbul
prayertime-cli times get --query Istanbul --date 2026-03-07 --json
prayertime-cli times get --lat 41.01384 --lon 28.94966 --field iftar --quiet
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

			return writeTimesHuman(cmd.OutOrStdout(), response, field, quiet)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve before fetching prayer times")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&date, "date", "today", "Date in YYYY-MM-DD format or 'today'")
	cmd.Flags().StringVar(&field, "field", "", "Return a single field such as maghrib, iftar, imsak, timezone, or method_name")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Emit only the selected field value")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate")
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
		Example: strings.TrimSpace(`
prayertime-cli times countdown --query Istanbul --target next-prayer
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

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve before fetching prayer times")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&target, "target", "", "Target prayer: next-prayer, imsak, fajr, sunrise, dhuhr, asr, maghrib, sunset, isha, iftar")
	cmd.Flags().StringVar(&at, "at", "", "Optional RFC3339 timestamp to evaluate countdown from")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Emit only remaining seconds")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		latitudeSet = cmd.Flags().Changed("lat")
		longitudeSet = cmd.Flags().Changed("lon")
	}

	return cmd
}
