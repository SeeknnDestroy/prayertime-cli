package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/SeeknnDestroy/prayertime-cli/internal/providers/aladhan"
	"github.com/SeeknnDestroy/prayertime-cli/internal/providers/openmeteo"
	"github.com/SeeknnDestroy/prayertime-cli/internal/version"
	"github.com/spf13/cobra"
)

type Dependencies struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Resolver app.LocationResolver
	Provider app.PrayerTimeProvider
	Clock    app.Clock
}

func Execute() int {
	deps := defaultDependencies()
	cmd := NewRootCmd(deps)

	if err := cmd.Execute(); err != nil {
		return renderCommandError(deps.Stdout, deps.Stderr, isJSONEnabled(cmd), err)
	}

	return app.ExitSuccess
}

func NewRootCmd(deps Dependencies) *cobra.Command {
	if deps.Stdout == nil {
		deps.Stdout = os.Stdout
	}
	if deps.Stderr == nil {
		deps.Stderr = os.Stderr
	}
	if deps.Clock == nil {
		deps.Clock = app.SystemClock{}
	}
	if deps.Resolver == nil || deps.Provider == nil {
		defaults := defaultDependencies()
		if deps.Resolver == nil {
			deps.Resolver = defaults.Resolver
		}
		if deps.Provider == nil {
			deps.Provider = defaults.Provider
		}
	}

	service := app.NewService(deps.Resolver, deps.Provider, deps.Clock)
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "prayertime-cli",
		Short: "CLI-first, agent-native Islamic prayer times tool",
		Long: strings.TrimSpace(`
CLI-first, agent-native Islamic prayer times tool.

Contract:
  - Structured payloads go to stdout.
  - Human-readable errors and suggestions go to stderr.
  - With --output json, errors are emitted as JSON on stdout.

Global output modes:
  - --output text prints human-readable output.
  - --output json prints structured JSON.
  - --output value is reserved for commands that expose --field selectors.

Exit codes:
  - 0 success
  - 1 internal failure
  - 2 usage error
  - 3 not found or ambiguous input
  - 4 network or upstream timeout
  - 5 reserved conflict/state error
`),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetOut(deps.Stdout)
	cmd.SetErr(deps.Stderr)
	cmd.PersistentFlags().String("output", string(outputText), "Output mode: text, json, or value")
	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Compatibility alias for --output json")
	_ = cmd.PersistentFlags().MarkHidden("json")
	cmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		return app.NewUsageError(err.Error(), "", fmt.Sprintf("Run '%s --help' for usage.", command.CommandPath()))
	})
	cmd.Version = version.String()
	cmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}")
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newLocationsCmd(service))
	cmd.AddCommand(newTimesCmd(service))

	return cmd
}

func isJSONEnabled(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}

	mode, err := resolveOutputMode(cmd, false)
	return err == nil && mode == outputJSON
}

func defaultDependencies() Dependencies {
	httpClient := &http.Client{Timeout: 8 * time.Second}
	return Dependencies{
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Resolver: openmeteo.NewClient(httpClient, ""),
		Provider: aladhan.NewClient(httpClient, ""),
		Clock:    app.SystemClock{},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := resolveOutputMode(cmd, false)
			if err != nil {
				return err
			}
			if mode == outputJSON {
				return writeJSON(cmd.OutOrStdout(), map[string]string{
					"version": version.Tag(),
					"commit":  version.Commit(),
					"date":    version.Date(),
				})
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), version.String())
			return err
		},
	}
}
