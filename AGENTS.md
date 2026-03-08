# AGENTS.md

## Project Purpose

- `prayertime-cli` is a CLI-first, agent-native Islamic prayer times tool.
- MVP 1 is stateless and limited to location search, daily prayer time lookup, and countdown queries.

## Source Of Truth

- Resolve place names with Open-Meteo geocoding.
- Fetch prayer times from AlAdhan with `method=13`.
- Treat coordinates as canonical once a place is resolved.

## CLI Contract

- Commands and flags are English-first.
- Turkish support is limited to semantic aliases for prayer identifiers and field selectors.
- MVP 1 has no persisted default location.
- Every `times` command requires `--query <place>` or both `--lat` and `--lon`.
- `--date today` is evaluated in the resolved location timezone.
- Structured payloads go to `stdout`; human-readable errors and suggestions go to `stderr`.
- With `--json`, error payloads are emitted as JSON on `stdout`.
- Preserve the exit-code contract:
  - `0`: success
  - `1`: internal failure
  - `2`: usage error
  - `3`: not found or ambiguous input
  - `4`: network or upstream timeout
  - `5`: reserved conflict/state error
- Use the compiled binary for exit-code examples. `go run` wraps non-zero process exits.

## Documentation Layout

- Hand-written, high-signal docs:
  - `README.md`
  - `README.tr.md`
  - `docs/agent-workflows.md`
  - `docs/cli-contract.md`
  - `docs/agent-evaluation.md`
  - `llms.txt`
- Generated docs:
  - `docs/cli/`
  - `docs/man/`
  - `completions/`
  - `llms-full.txt`
- Do not hand-edit generated docs. Regenerate them with the docs command.

## Commands

- Test: `go test ./...`
- Build: `go build ./cmd/prayertime-cli`
- Regenerate docs: `go run ./cmd/prayertime-cli-docs`
- Verify generated docs: `make docs-check`

## Repository Workflow

- Work on feature branches only.
- Keep commits atomic.
- Prefer stacked PRs when changes separate cleanly.
