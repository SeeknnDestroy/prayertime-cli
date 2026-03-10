# prayertime-cli

`prayertime-cli` is a stateless CLI for Islamic prayer schedules and countdowns. It is built for agents, shell scripts, and direct terminal use.

## MVP 1

- Search locations with Open-Meteo geocoding.
- Fetch daily prayer times from AlAdhan with `method=13` (Diyanet).
- Count down to the next prayer or a named prayer.
- Emit structured JSON on `stdout` with `--json` or bare scalar values with `--quiet`.

## Input Model

- MVP 1 has no persisted default location.
- Every `times` command requires either `--query <place>` or both `--lat <float>` and `--lon <float>`.
- `--date today` is resolved in the target location timezone.

## Common Tasks

Find candidate locations before choosing one:

```bash
prayertime-cli locations search --query "Springfield" --country-code US --json
```

Get today's full prayer schedule:

```bash
prayertime-cli times get --query Istanbul --json
```

Extract one value for automation or shell pipelines:

```bash
prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
```

Ask how long until the next ezan / next prayer:

```bash
prayertime-cli times countdown --query Istanbul --target next-prayer --json
```

Ask how long until a specific prayer, including iftar:

```bash
prayertime-cli times countdown --query Istanbul --target iftar --quiet
```

Use coordinates instead of a place name:

```bash
prayertime-cli times get --lat 41.01384 --lon 28.94966 --date today --json
```

Recover from a typo or ambiguous location:

```bash
prayertime-cli locations search --query Istnbul --json
```

## Output Modes

- `--json`: emit structured payloads to `stdout`. With `--json`, errors are also JSON on `stdout`.
- `--quiet`: emit one bare scalar value. `times get` requires `--field`; `times countdown --quiet` defaults to `seconds_remaining`.
- `--output text|json|value`: generalized form of the same output contract. Use `--output` when you want one explicit output switch across commands.
- Default human mode: readable output on `stdout`; errors and suggestions on `stderr`.
- If you need exact process exit codes, run the compiled binary directly. `go run` wraps non-zero exits.

## Aliases

- Commands and flags are English-first and canonical.
- Turkish semantic aliases are supported for prayer targets and field selectors.
- `iftar` and `aksam` resolve to `maghrib`; `yatsi` resolves to `isha`; `ogle` resolves to `dhuhr`.

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

## Build And Docs

```bash
make verify
make docs
make build
make release-check
```

## Exit Codes

- `0`: success
- `1`: internal failure
- `2`: usage error
- `3`: not found or ambiguous input
- `4`: network or upstream timeout
- `5`: reserved conflict/state error

## More Docs

- [Agent Workflows](docs/agent-workflows.md)
- [CLI Contract](docs/cli-contract.md)
- [Agent Evaluation](docs/agent-evaluation.md)
- [CLI Reference](docs/cli/prayertime-cli.md)
- [ADR 0002: Data Sources](docs/adr/0002-data-sources.md)

## License

MIT
