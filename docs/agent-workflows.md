# Agent Workflows

> Natural-language intents mapped to stable `prayertime-cli` commands for MVP 1.

## Working Model

- MVP 1 is stateless. There is no saved default location.
- Every `times` command needs `--query <place>` or both `--lat` and `--lon`.
- Use `locations search` first when a place may be ambiguous or misspelled.

## Resolve A Place

Use this when the user asks for a city that may have multiple matches or an uncertain spelling.

```bash
prayertime-cli locations search --query "Springfield" --country-code US --limit 3 --json
prayertime-cli locations search --query Istnbul --json
```

Good intents:

- "Which Springfield do you mean?"
- "Search for Konya, Türkiye"
- "I typed Istnbul, show likely matches"

## Get Today's Full Schedule

Use `times get` when the user wants the full daily prayer table.

```bash
prayertime-cli times get --query Istanbul --json
prayertime-cli times get --lat 41.01384 --lon 28.94966 --date today --json
```

Good intents:

- "Bugün namaz vakitleri"
- "Show today's prayer times in Ankara, Türkiye"
- "Give me today's schedule for these coordinates"

## Extract One Field

Use `--field` with scalar output when the caller wants one specific value for synthesis, piping, or follow-on automation.

```bash
prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
prayertime-cli times get --query Istanbul --field timezone --quiet
prayertime-cli times get --query Istanbul --field iftar --quiet
```

Good intents:

- "Yalnızca yatsı vaktini ver"
- "What timezone is this schedule in?"
- "Give me today's iftar time"

## Generic Next-Prayer Countdown

Use `next-prayer` for broad countdown questions like "how long until the next ezan?" With no `--field`, scalar countdown output defaults to `minutes_remaining`.

```bash
prayertime-cli times countdown --query Istanbul --target next-prayer --json
prayertime-cli times countdown --lat 41.01384 --lon 28.94966 --target next-prayer --quiet
```

Good intents:

- "Ezana kaç dakika kaldı?"
- "How long until the next prayer?"
- "Give me minutes until the next prayer"

## Specific-Prayer Countdown

Use named targets when the user asks for a specific prayer. Canonical targets are English; Turkish aliases are accepted. With no `--field`, scalar countdown output defaults to `minutes_remaining`.

```bash
prayertime-cli times countdown --query Istanbul --target asr --json
prayertime-cli times countdown --query Istanbul --target iftar --quiet
prayertime-cli times countdown --query Istanbul --target yatsi --json
```

Good intents:

- "How long until asr?"
- "İftara kaç dakika kaldı?"
- "Yatsıya ne kadar kaldı?"

## Evaluate From A Specific Time

Use `--at` for deterministic replay, testing, or evaluation runs.

```bash
prayertime-cli times countdown --query Istanbul --target next-prayer --at 2026-03-07T18:00:00+03:00 --json
```

Good intents:

- "What would the countdown have been at 18:00?"
- "Replay this workflow against a fixed timestamp"

## Recovery Patterns

- Missing location input:
  - retry with `--query <place>` or both `--lat` and `--lon`
- Not found:
  - run `locations search --query <text> --json`
- Ambiguous place:
  - run `locations search --query <text> --json` and pick coordinates or a more exact place
- Network failure:
  - retry with backoff; keep the same arguments

## Output Mode Selection

- Use `--json` for structured payloads and machine parsing.
- Use `--quiet` when a command should emit a single scalar value.
- For countdown, `--quiet` and bare `--output value` default to `minutes_remaining`; use `--field` when you need a different scalar.
- Use `--output text|json|value` when an agent or wrapper wants one explicit output switch across commands.
- Use default human mode for local terminal inspection.
