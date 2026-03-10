package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/SeeknnDestroy/prayertime-cli/internal/cli"
	cobradoc "github.com/spf13/cobra/doc"
)

func main() {
	docsDir := flag.String("docs-dir", filepath.Join("docs", "cli"), "Directory for generated Markdown CLI docs")
	manDir := flag.String("man-dir", filepath.Join("docs", "man"), "Directory for generated man pages")
	completionsDir := flag.String("completions-dir", "completions", "Directory for generated shell completions")
	llmsFullPath := flag.String("llms-full-path", "llms-full.txt", "Output path for generated llms-full.txt")
	flag.Parse()

	root := cli.NewRootCmd(cli.Dependencies{
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	root.DisableAutoGenTag = true

	must(os.MkdirAll(*docsDir, 0o755))
	must(os.MkdirAll(*manDir, 0o755))
	must(os.MkdirAll(*completionsDir, 0o755))
	must(os.MkdirAll(filepath.Dir(*llmsFullPath), 0o755))

	must(cobradoc.GenMarkdownTree(root, *docsDir))
	header := &cobradoc.GenManHeader{
		Title:   "prayertime-cli",
		Section: "1",
		Source:  "prayertime-cli",
	}
	must(cobradoc.GenManTree(root, header, *manDir))
	must(root.GenBashCompletionFileV2(filepath.Join(*completionsDir, "prayertime-cli.bash"), true))
	must(root.GenZshCompletionFile(filepath.Join(*completionsDir, "_prayertime-cli")))
	must(root.GenFishCompletionFile(filepath.Join(*completionsDir, "prayertime-cli.fish"), true))
	must(root.GenPowerShellCompletionFile(filepath.Join(*completionsDir, "prayertime-cli.ps1")))
	must(writeLLMSFull(*llmsFullPath))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
