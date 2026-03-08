## prayertime-cli times get

Fetch prayer times for a location and date

### Synopsis

Fetch one day's prayer schedule for a resolved place or explicit coordinates.

MVP 1 is stateless, so every call must include --query PLACE or both --lat and --lon. Use --field with --quiet for one scalar value, or --json for the full structured response.

```
prayertime-cli times get [flags]
```

### Examples

```
prayertime-cli times get --query Istanbul
prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
prayertime-cli times get --lat 41.01384 --lon 28.94966 --date 2026-03-07 --json
```

### Options

```
      --country-code string   Optional ISO country code filter
      --date string           Date in YYYY-MM-DD format or 'today' (default "today")
      --field string          Return one field such as maghrib, iftar, yatsi, timezone, or method_name
  -h, --help                  help for get
      --lat float             Latitude coordinate. Use with --lon instead of --query
      --lon float             Longitude coordinate. Use with --lat instead of --query
      --query string          Place name to resolve. Required unless --lat and --lon are set
  -q, --quiet                 Emit only the selected field value
```

### Options inherited from parent commands

```
      --json   Emit structured JSON to stdout
```

### SEE ALSO

* [prayertime-cli times](prayertime-cli_times.md)	 - Fetch daily prayer schedules and countdowns

