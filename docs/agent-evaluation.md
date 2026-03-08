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
- `--json` returns structured payloads and JSON error objects
- missing location input returns exit code `2`
- not-found or ambiguous input returns exit code `3`
