package main

import (
	"fmt"
	"os"
	"strings"
)

type llmsSource struct {
	Path        string
	Description string
}

var llmsSources = []llmsSource{
	{
		Path:        "README.md",
		Description: "Fast entry point covering the CLI purpose, common tasks, output modes, aliases, and exit codes.",
	},
	{
		Path:        "AGENTS.md",
		Description: "Compact repository briefing for coding agents, including source-of-truth rules and exact build, test, and docs commands.",
	},
	{
		Path:        "docs/agent-workflows.md",
		Description: "Task-first workflow guide mapping natural-language intents to stable commands.",
	},
	{
		Path:        "docs/cli-contract.md",
		Description: "Stable CLI contract for inputs, outputs, aliases, error payloads, and exit codes.",
	},
	{
		Path:        "docs/adr/0001-go-stack.md",
		Description: "Stack decision for Go and Cobra.",
	},
	{
		Path:        "docs/adr/0002-data-sources.md",
		Description: "Source-of-truth decision for Open-Meteo and AlAdhan method 13.",
	},
}

func writeLLMSFull(path string) error {
	var builder strings.Builder
	builder.WriteString("# prayertime-cli\n\n")
	builder.WriteString("> Consolidated, high-signal documentation for agent ingestion. Generated from a fixed source list by `go run ./cmd/prayertime-cli-docs`.\n\n")
	builder.WriteString("## Included Sources\n\n")

	for _, source := range llmsSources {
		_, _ = fmt.Fprintf(&builder, "- `%s`: %s\n", source.Path, source.Description)
	}

	builder.WriteString("\n")

	for _, source := range llmsSources {
		content, err := os.ReadFile(source.Path)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(&builder, "## Source: `%s`\n\n", source.Path)
		builder.Write(content)
		if !strings.HasSuffix(builder.String(), "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	return os.WriteFile(path, []byte(builder.String()), 0o644)
}
