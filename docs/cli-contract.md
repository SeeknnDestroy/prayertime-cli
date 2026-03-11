# CLI Contract

> Stable agent-facing behavior for `prayertime-cli` MVP 1.

## Scope

- `locations search` resolves place names to candidate coordinates.
- `times get` returns one day's prayer schedule.
- `times countdown` returns the next-prayer or named-prayer countdown.
- MVP 1 is stateless. There is no saved location or config fallback.

## Required Inputs

- `times get` and `times countdown` require one location selector:
  - `--query <place>`
  - or `--lat <float> --lon <float>`
- `--query` and `--lat/--lon` are mutually exclusive.
- `--date today` is evaluated in the resolved location timezone.
- `times countdown --target next-prayer` is the generic "next ezan" path.

## Output Modes

- Default mode:
  - success payloads are human-readable text on `stdout`
  - errors and suggestions are written to `stderr`
- `--json`:
  - success payloads are JSON on `stdout`
  - errors are JSON on `stdout`
- `--quiet`:
  - `times get` requires `--field`
  - `times get --quiet` emits only the selected field value
  - `times countdown --quiet` emits `minutes_remaining` when `--field` is omitted
- `--output text|json|value`:
  - generalized form of the same output model
  - `--output json` is equivalent to `--json`
  - `times countdown --output value` defaults to `minutes_remaining` when `--field` is omitted
  - otherwise `--output value` requires a command-specific scalar selector such as `--field`

## Alias Rules

- Commands and flags stay English-first.
- Turkish semantic aliases are accepted for prayer targets and field selectors.
- Common target aliases:
  - `iftar`, `aksam` -> `maghrib`
  - `yatsi` -> `isha`
  - `ogle` -> `dhuhr`
  - `gunes` -> `sunrise`
- Common field aliases follow the same prayer mapping and resolve to `_at` fields, for example:
  - `iftar` -> `maghrib_at`
  - `yatsi` -> `isha_at`

## Representative Success Payloads

`times get --query Istanbul --json`

```json
{
  "location_name": "Istanbul, Türkiye",
  "latitude": 41.01384,
  "longitude": 28.94966,
  "timezone": "Europe/Istanbul",
  "date": "2026-03-09",
  "imsak_at": "2026-03-09T05:45:00+03:00",
  "fajr_at": "2026-03-09T05:55:00+03:00",
  "sunrise_at": "2026-03-09T07:19:00+03:00",
  "dhuhr_at": "2026-03-09T13:20:00+03:00",
  "asr_at": "2026-03-09T16:33:00+03:00",
  "maghrib_at": "2026-03-09T19:11:00+03:00",
  "sunset_at": "2026-03-09T19:11:00+03:00",
  "isha_at": "2026-03-09T20:30:00+03:00",
  "method_id": 13,
  "method_name": "Diyanet İşleri Başkanlığı, Turkey (experimental)",
  "source": "aladhan:method=13",
  "ramadan_active": true
}
```

`times countdown --query Istanbul --target next-prayer --json`

```json
{
  "location_name": "Istanbul, Türkiye",
  "latitude": 41.01384,
  "longitude": 28.94966,
  "timezone": "Europe/Istanbul",
  "date": "2026-03-09",
  "imsak_at": "2026-03-09T05:45:00+03:00",
  "fajr_at": "2026-03-09T05:55:00+03:00",
  "sunrise_at": "2026-03-09T07:19:00+03:00",
  "dhuhr_at": "2026-03-09T13:20:00+03:00",
  "asr_at": "2026-03-09T16:33:00+03:00",
  "maghrib_at": "2026-03-09T19:11:00+03:00",
  "sunset_at": "2026-03-09T19:11:00+03:00",
  "isha_at": "2026-03-09T20:30:00+03:00",
  "method_id": 13,
  "method_name": "Diyanet İşleri Başkanlığı, Turkey (experimental)",
  "source": "aladhan:method=13",
  "ramadan_active": true,
  "target": "maghrib",
  "target_at": "2026-03-09T19:11:00+03:00",
  "seconds_remaining": 65798,
  "minutes_remaining": 1096
}
```

## Error Payload

`times countdown --query Istnbul --target iftar --json`

```json
{
  "ok": false,
  "exit_code": 3,
  "error_type": "not_found",
  "message": "no locations matched \"Istnbul\"",
  "input_received": "Istnbul",
  "suggestion": "Run 'prayertime-cli locations search --query \"Istnbul\" --json' to inspect candidates."
}
```

## Exit Codes

- `0`: success
- `1`: internal failure
- `2`: usage error
- `3`: not found or ambiguous input
- `4`: network or upstream timeout
- `5`: reserved conflict/state error

Run the compiled binary directly when you need exact exit-code behavior. `go run` wraps failures and does not preserve the tool's process exit code.
