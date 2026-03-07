## prayertime-cli times get

Fetch prayer times for a location and date

```
prayertime-cli times get [flags]
```

### Examples

```
prayertime-cli times get --query Istanbul
prayertime-cli times get --query Istanbul --date 2026-03-07 --json
prayertime-cli times get --lat 41.01384 --lon 28.94966 --field iftar --quiet
```

### Options

```
      --country-code string   Optional ISO country code filter
      --date string           Date in YYYY-MM-DD format or 'today' (default "today")
      --field string          Return a single field such as maghrib, iftar, imsak, timezone, or method_name
  -h, --help                  help for get
      --lat float             Latitude coordinate
      --lon float             Longitude coordinate
      --query string          Place name to resolve before fetching prayer times
  -q, --quiet                 Emit only the selected field value
```

### Options inherited from parent commands

```
      --json   Emit JSON output to stdout
```

### SEE ALSO

* [prayertime-cli times](prayertime-cli_times.md)	 - Fetch daily prayer times and countdowns

