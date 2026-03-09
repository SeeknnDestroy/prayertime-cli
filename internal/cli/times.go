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
		Long: strings.TrimSpace(`
Fetch daily prayer times and countdowns for a resolved location.

Contract:
  - Provide exactly one location input strategy: --query <place> or --lat with --lon.
  - --output text prints human-readable output to stdout.
  - --output json prints structured JSON to stdout.
  - --output value prints only the selected --field value.
  - Human-readable errors go to stderr. With --output json, errors are JSON on stdout.
`),
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
	var view string
	var quiet bool
	var latitude float64
	var longitude float64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch prayer times for a location and date",
		Long: strings.TrimSpace(`
Resolve a place or accept explicit coordinates, then return daily prayer times.

Input rules:
  - Provide exactly one of --query <place> or --lat <value> --lon <value>.
  - --query may be paired with --country-code to narrow ambiguous locations.
  - --date accepts YYYY-MM-DD or today.
  - --date today is evaluated in the resolved location timezone.

Output:
  - --output text prints human-readable prayer times.
  - --output json prints structured JSON.
  - --output value prints only the selected --field value.
  - --field accepts canonical fields such as maghrib_at, timezone, method_name, and source.
  - With --field and --output json, the response is a single-key JSON object.

View modes:
  - concise: location_name, timezone, date, prayer times, ramadan_active
  - detailed: concise fields plus latitude, longitude, method_id, method_name, source
  - Default view is detailed for text and json output.
`),
		Example: strings.TrimSpace(`
prayertime-cli times get --query Istanbul
prayertime-cli times get --query Istanbul --date 2026-03-07 --output json
prayertime-cli times get --query Istanbul --field maghrib_at --output value
prayertime-cli times get --lat 41.01384 --lon 28.94966 --view concise --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := resolveOutputMode(cmd, quiet)
			if err != nil {
				return err
			}

			if mode == outputValue && field == "" {
				return app.NewUsageError(
					"--output value requires --field for times get",
					"",
					"Provide --field to choose a single prayer-times value.",
				).WithDetails(app.ErrorDetails{ValidFields: app.ValidTimesFields()})
			}
			if err := validateTimesField(field); err != nil {
				return err
			}

			resolvedView, err := resolveViewMode(view, viewDetailed)
			if err != nil {
				return err
			}

			request := app.TimesRequest{
				Query:       query,
				CountryCode: strings.ToUpper(countryCode),
				Date:        date,
			}
			if cmd.Flags().Changed("lat") {
				request.Latitude = &latitude
			}
			if cmd.Flags().Changed("lon") {
				request.Longitude = &longitude
			}

			response, err := service.GetTimes(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeTimesOutput(cmd.OutOrStdout(), response, mode, resolvedView, field)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve before fetching prayer times")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&date, "date", "today", "Date in YYYY-MM-DD format or 'today'")
	cmd.Flags().StringVar(&field, "field", "", "Single canonical field to return, such as maghrib_at or method_name")
	cmd.Flags().StringVar(&view, "view", "", "Response view: concise or detailed")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Compatibility alias for --output value")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate")
	cmd.MarkFlagsMutuallyExclusive("query", "lat")
	cmd.MarkFlagsMutuallyExclusive("query", "lon")
	cmd.MarkFlagsRequiredTogether("lat", "lon")
	_ = cmd.Flags().MarkHidden("quiet")

	return cmd
}

func newTimesCountdownCmd(service *app.Service) *cobra.Command {
	var query string
	var countryCode string
	var target string
	var at string
	var field string
	var view string
	var quiet bool
	var latitude float64
	var longitude float64

	cmd := &cobra.Command{
		Use:   "countdown",
		Short: "Calculate seconds and minutes remaining until a target prayer",
		Long: strings.TrimSpace(`
Resolve a place or accept explicit coordinates, then calculate the remaining time until a target prayer.

Input rules:
  - Provide exactly one of --query <place> or --lat <value> --lon <value>.
  - --target is required.
  - --target accepts canonical values plus Turkish aliases such as iftar, öğle, akşam, and yatsı.
  - --at accepts an RFC3339 timestamp. If omitted, the current time is used.

Output:
  - --output text prints a human-readable countdown.
  - --output json prints structured JSON.
  - --output value prints only the selected --field value.
  - --field accepts canonical countdown fields such as seconds_remaining, target_at, maghrib_at, and method_name.
  - With --field and --output json, the response is a single-key JSON object.
  - Legacy --quiet maps to --output value and defaults to seconds_remaining.

View modes:
  - concise: location_name, timezone, date, target, target_at, seconds_remaining, minutes_remaining
  - detailed: concise fields plus the full detailed prayer-times payload
  - Default view is concise for text and json output.
`),
		Example: strings.TrimSpace(`
prayertime-cli times countdown --query Istanbul --target next-prayer
prayertime-cli times countdown --query Istanbul --target iftar --field seconds_remaining --output value
prayertime-cli times countdown --lat 41.01384 --lon 28.94966 --target asr --at 2026-03-07T12:00:00+03:00 --output json
prayertime-cli times countdown --query Istanbul --target maghrib --view detailed --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := resolveOutputMode(cmd, quiet)
			if err != nil {
				return err
			}

			resolvedField := field
			if mode == outputValue && resolvedField == "" {
				if quiet {
					resolvedField = "seconds_remaining"
				} else {
					return app.NewUsageError(
						"--output value requires --field for times countdown",
						"",
						"Provide --field to choose a single countdown value.",
					).WithDetails(app.ErrorDetails{ValidFields: app.ValidCountdownFields()})
				}
			}
			if err := validateCountdownField(resolvedField); err != nil {
				return err
			}

			resolvedView, err := resolveViewMode(view, viewConcise)
			if err != nil {
				return err
			}

			var atValue *time.Time
			if strings.TrimSpace(at) != "" {
				parsed, err := time.Parse(time.RFC3339, at)
				if err != nil {
					return app.NewUsageError(
						fmt.Sprintf("invalid RFC3339 value %q", at),
						at,
						"Use --at 2026-03-07T18:00:00+03:00.",
					)
				}
				atValue = &parsed
			}

			request := app.CountdownRequest{
				Query:       query,
				CountryCode: strings.ToUpper(countryCode),
				Target:      target,
				At:          atValue,
			}
			if cmd.Flags().Changed("lat") {
				request.Latitude = &latitude
			}
			if cmd.Flags().Changed("lon") {
				request.Longitude = &longitude
			}

			response, err := service.GetCountdown(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeCountdownOutput(cmd.OutOrStdout(), response, mode, resolvedView, resolvedField)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to resolve before fetching prayer times")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().StringVar(&target, "target", "", "Target prayer such as next-prayer, imsak, fajr, sunrise, dhuhr, asr, maghrib, sunset, isha, or iftar")
	cmd.Flags().StringVar(&at, "at", "", "Optional RFC3339 timestamp to evaluate countdown from")
	cmd.Flags().StringVar(&field, "field", "", "Single canonical field to return, such as seconds_remaining or target_at")
	cmd.Flags().StringVar(&view, "view", "", "Response view: concise or detailed")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Compatibility alias for --output value")
	cmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude coordinate")
	cmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude coordinate")
	_ = cmd.MarkFlagRequired("target")
	cmd.MarkFlagsMutuallyExclusive("query", "lat")
	cmd.MarkFlagsMutuallyExclusive("query", "lon")
	cmd.MarkFlagsRequiredTogether("lat", "lon")
	_ = cmd.Flags().MarkHidden("quiet")

	return cmd
}
