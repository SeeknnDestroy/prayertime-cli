package cli

import (
	"fmt"
	"os"

	"github.com/SeeknnDestroy/prayertime-cli/internal/version"
	"github.com/spf13/cobra"
)

func Execute() int {
	root := newRootCmd()

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "prayertime-cli",
		Short:         "CLI-first, agent-native Islamic prayer times tool",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Version = version.String()
	cmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}")
	cmd.AddCommand(newVersionCmd())

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version.String())
		},
	}
}
