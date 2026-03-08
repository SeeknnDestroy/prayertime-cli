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
Search for a place before requesting prayer times.

Use this command first when a user-supplied place may be ambiguous, incomplete, or misspelled.
`),
		Example: strings.TrimSpace(`
prayertime-cli locations search --query Istanbul
prayertime-cli locations search --query Springfield --country-code US --limit 3 --json
prayertime-cli locations search --query Istnbul --json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := service.SearchLocations(cmd.Context(), query, strings.ToUpper(countryCode), limit)
			if err != nil {
				return err
			}

			if isJSONEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), response)
			}

			return writeLocationsHuman(cmd.OutOrStdout(), response)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Place name to search")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Optional ISO country code filter")
	cmd.Flags().IntVar(&limit, "limit", 5, "Maximum number of results to return")

	return cmd
}

func writeLocationsHuman(out interface{ Write([]byte) (int, error) }, response app.LocationSearchResponse) error {
	for index, location := range response.Results {
		if _, err := fmt.Fprintf(out, "%d. %s [%.6f, %.6f] %s\n", index+1, location.DisplayName(), location.Latitude, location.Longitude, location.Timezone); err != nil {
			return err
		}
	}
	return nil
}
