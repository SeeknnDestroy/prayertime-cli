# Agent Evaluation

> Lightweight prompt and command checklist for validating the docs and help surface against realistic agent tasks.

## Prompt Set

- "Bugün namaz vakitleri Istanbul"
- "Ezana kaç dakika kaldı?"
- "How long until the next prayer in Istanbul?"
- "İftara kaç saniye kaldı?"
- "Yalnızca yatsı vaktini ver"
- "Use these coordinates and return today's full schedule"
- "Istnbul için doğru şehri bul"
- "Springfield belirsiz, seçenekleri göster"
- "Run the countdown from a fixed RFC3339 timestamp"
- "Return JSON and preserve the exit code on ambiguous input"

## Agent CLI Checklist

- Structured output is discoverable:
  - `--json` is visible in help and returns JSON on `stdout`
  - `--quiet` is visible on scalar-capable commands and emits a single bare value
  - `--output text|json|value` remains available for wrappers that want one explicit switch
- Errors are automation-friendly:
  - human-readable errors go to `stderr`
  - `--json` returns structured error payloads on `stdout`
  - exit codes stay stable and documented
- Commands are safe and predictable:
  - MVP 1 is stateless and read-only, so query operations are naturally idempotent
  - no destructive commands exist, so `--dry-run` and `--yes` are intentionally not applicable in this version
- Help is self-documenting:
  - `--help` must expose examples for `--json`, `--quiet`, field selection, and deterministic countdown replay
- Composability is covered:
  - `--field` plus `--quiet` or `--output value` supports scalar pipelines
  - countdown scalar defaults resolve to `minutes_remaining` when no `--field` is provided
  - full JSON payloads stay stable for agent parsing

## Verification Checklist

Use the compiled binary for exit-code checks:

```bash
go build -o ./bin/prayertime-cli ./cmd/prayertime-cli
```

Run these checks:

```bash
./bin/prayertime-cli locations search --query Istnbul --json
./bin/prayertime-cli locations search --query Springfield --country-code US --limit 3 --json
./bin/prayertime-cli times get --query Istanbul --json
./bin/prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
./bin/prayertime-cli times get --lat 41.01384 --lon 28.94966 --date today --json
./bin/prayertime-cli times countdown --query Istanbul --target next-prayer --json
./bin/prayertime-cli times countdown --query Istanbul --target iftar --quiet
./bin/prayertime-cli times countdown --query Istanbul --target next-prayer --at 2026-03-07T18:00:00+03:00 --json
./bin/prayertime-cli times countdown --json
```

Confirm:

- task-first docs point to the same commands shown in `--help`
- `next-prayer` is the generic countdown path
- `iftar` works as a supported alias, not the default framing for all countdowns
- `--quiet` returns scalar-only output
- `times countdown --quiet` and bare `--output value` default to `minutes_remaining`
- `--json` returns structured payloads and JSON error objects
- `--output json` and `--output value` remain equivalent to the shortcut flags
- missing location input returns exit code `2`
- not-found or ambiguous input returns exit code `3`
