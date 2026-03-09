package cli

import (
	"fmt"
	"strings"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/spf13/cobra"
)

func newLocationsCmd(service *app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locations",
		Short: "Search locations before requesting prayer times",
		Long: strings.TrimSpace(`
Search and inspect candidate places before requesting prayer times.

Contract:
  - --query is required.
  - Use --country-code to narrow ambiguous place names.
  - --output text prints numbered candidates for humans.
  - --output json prints structured candidates with display_name, coordinates, and timezone.
  - Success payloads go to stdout. Errors go to stderr unless --output json is used.
`),
	}

	cmd.AddCommand(newLocationsSearchCmd(service))
	return cmd
}

func newLocationsSearchCmd(service *app.Service) *cobra.Command {
	var query string
	var countryCode string
	var limit int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for a place and return candidate coordinates",
		Long: strings.TrimSpace(`
Resolve a place query into candidate locations and canonical coordinates.

Input rules:
  - --query is required.
  - --country-code is optional and should be used to narrow ambiguous place names.
  - --limit controls the maximum number of returned candidates.

Output:
  - --output text prints numbered candidates with coordinates and timezone.
  - --output json prints query, count, and structured candidates including display_name.
  - --output value is not supported for location search.
`),
		Example: strings.TrimSpace(`
prayertime-cli locations search --query Istanbul
prayertime-cli locations search --query London --country-code GB --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := resolveOutputMode(cmd, false)
			if err != nil {
				return err
			}
			if mode == outputValue {
				return app.NewUsageError(
					"--output value is not supported for locations search",
					string(mode),
					"Use --output text or --output json.",
				)
			}

			response, err := service.SearchLocations(cmd.Context(), query, strings.ToUpper(countryCode), limit)
			if err != nil {
				return err
			}

			if mode == outputJSON {
				return writeJSON(cmd.OutOrStdout(), response)
			}

			return writeLocationsHuman(cmd.OutOrStdout(), response)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to search")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().IntVar(&limit, "limit", 5, "Maximum number of results to return")
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

func writeLocationsHuman(out interface{ Write([]byte) (int, error) }, response app.LocationSearchResponse) error {
	for index, location := range response.Results {
		if _, err := fmt.Fprintf(out, "%d. %s [%.6f, %.6f] %s\n", index+1, location.DisplayName, location.Latitude, location.Longitude, location.Timezone); err != nil {
			return err
		}
	}
	return nil
}
