# AGENTS.md

## Project Purpose

- `prayertime-cli` is a CLI-first, agent-native Islamic prayer times tool.
- MVP 1 is stateless and limited to location search, daily prayer time lookup, and countdown queries.

## Implementation Rules

- Keep commands and flags English-first. Turkish support is limited to semantic aliases for prayer identifiers and field selectors.
- Treat Open-Meteo as the location resolver and AlAdhan `method=13` as the prayer time source of truth for MVP 1.
- Keep structured payloads on `stdout` and diagnostics on `stderr`.
- Maintain the documented exit-code contract:
  - `0`: success
  - `1`: internal failure
  - `2`: usage error
  - `3`: not found or ambiguous input
  - `4`: network or upstream timeout
  - `5`: reserved conflict/state error

## Repository Workflow

- Work on feature branches only and keep commits atomic.
- Prefer stacked PRs for phased delivery when changes are logically separable.
- Run `go test ./...` before opening or updating a PR.

