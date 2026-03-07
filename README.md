# prayertime-cli

`prayertime-cli` is a CLI-first, agent-native Islamic prayer times tool built for deterministic automation, fast terminal workflows, and privacy-respecting daily use.

## Goals

- Provide exact daily prayer times and countdowns through a stable CLI contract.
- Optimize for AI agents and shell scripts with strict JSON support and predictable exit codes.
- Start with a stateless MVP powered by Open-Meteo geocoding and AlAdhan method `13` (Diyanet).

## Planned Command Surface

```text
prayertime-cli locations search --query <text> [--country-code TR] [--limit 5] [--json]
prayertime-cli times get (--query <text> | --lat <float> --lon <float>) [--country-code TR] [--date YYYY-MM-DD|today] [--json] [--field <key>] [--quiet]
prayertime-cli times countdown (--query <text> | --lat <float> --lon <float>) --target next-prayer|fajr|sunrise|dhuhr|asr|maghrib|isha|imsak|iftar [--at RFC3339] [--json] [--quiet]
prayertime-cli version
```

## Principles

- English commands and flags are canonical.
- Turkish semantic aliases are accepted only for prayer identifiers and field selectors.
- JSON responses are emitted to `stdout`; diagnostics stay on `stderr`.
- The CLI never prompts interactively.

## Development

This repository uses Go 1.26 and a phased, stacked-PR workflow.

```bash
go test ./...
go build ./cmd/prayertime-cli
```

## Install

Tagged releases are published as cross-platform binaries. Package manager automation is wired for Homebrew Cask and Scoop.

```bash
# Homebrew
brew tap SeeknnDestroy/homebrew-tap
brew install --cask prayertime-cli

# Scoop
scoop bucket add prayertime-cli https://github.com/SeeknnDestroy/scoop-bucket
scoop install prayertime-cli

# Go
go install github.com/SeeknnDestroy/prayertime-cli/cmd/prayertime-cli@latest
```

## Examples

```bash
prayertime-cli locations search --query Istanbul --country-code TR --json
prayertime-cli times get --query Istanbul --field iftar --quiet
prayertime-cli times countdown --query Istanbul --target next-prayer --json
```

## Exit Codes

- `0`: success
- `1`: internal failure
- `2`: usage error
- `3`: not found or ambiguous input
- `4`: network or upstream timeout
- `5`: reserved conflict/state error

## License

MIT
