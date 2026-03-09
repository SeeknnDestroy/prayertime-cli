package cli

import (
	"fmt"
	"strings"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/spf13/cobra"
)

type outputMode string

const (
	outputText  outputMode = "text"
	outputJSON  outputMode = "json"
	outputValue outputMode = "value"
)

type viewMode string

const (
	viewConcise  viewMode = "concise"
	viewDetailed viewMode = "detailed"
)

func resolveOutputMode(cmd *cobra.Command, quietAlias bool) (outputMode, error) {
	outputValueRaw, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", err
	}

	mode := outputMode(strings.TrimSpace(strings.ToLower(outputValueRaw)))
	if mode == "" {
		mode = outputText
	}

	switch mode {
	case outputText, outputJSON, outputValue:
	default:
		return "", app.NewUsageError(
			fmt.Sprintf("unsupported output mode %q", outputValueRaw),
			outputValueRaw,
			"Use --output text, --output json, or --output value.",
		)
	}

	jsonAlias := cmd.Flags().Changed("json")
	if jsonAlias {
		if mode != outputText && mode != outputJSON {
			return "", app.NewUsageError(
				"--json cannot be combined with a different --output mode",
				string(mode),
				"Use either --json or --output json.",
			)
		}
		mode = outputJSON
	}

	if quietAlias {
		if mode != outputText && mode != outputValue {
			return "", app.NewUsageError(
				"--quiet cannot be combined with a different --output mode",
				string(mode),
				"Use either --quiet or --output value.",
			)
		}
		mode = outputValue
	}

	return mode, nil
}

func resolveViewMode(input string, defaultMode viewMode) (viewMode, error) {
	trimmed := strings.TrimSpace(strings.ToLower(input))
	if trimmed == "" {
		return defaultMode, nil
	}

	mode := viewMode(trimmed)
	switch mode {
	case viewConcise, viewDetailed:
		return mode, nil
	default:
		return "", app.NewUsageError(
			fmt.Sprintf("unsupported view %q", input),
			input,
			"Use --view concise or --view detailed.",
		)
	}
}
