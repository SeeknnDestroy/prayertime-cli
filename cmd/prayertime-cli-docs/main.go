package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/SeeknnDestroy/prayertime-cli/internal/cli"
	cobradoc "github.com/spf13/cobra/doc"
)

func main() {
	root := cli.NewRootCmd(cli.Dependencies{
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	root.DisableAutoGenTag = true

	must(os.MkdirAll(filepath.Join("docs", "cli"), 0o755))
	must(os.MkdirAll(filepath.Join("docs", "man"), 0o755))
	must(os.MkdirAll("completions", 0o755))

	must(cobradoc.GenMarkdownTree(root, filepath.Join("docs", "cli")))
	header := &cobradoc.GenManHeader{
		Title:   "prayertime-cli",
		Section: "1",
		Source:  "prayertime-cli",
	}
	must(cobradoc.GenManTree(root, header, filepath.Join("docs", "man")))
	must(root.GenBashCompletionFileV2(filepath.Join("completions", "prayertime-cli.bash"), true))
	must(root.GenZshCompletionFile(filepath.Join("completions", "_prayertime-cli")))
	must(root.GenFishCompletionFile(filepath.Join("completions", "prayertime-cli.fish"), true))
	must(root.GenPowerShellCompletionFile(filepath.Join("completions", "prayertime-cli.ps1")))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
