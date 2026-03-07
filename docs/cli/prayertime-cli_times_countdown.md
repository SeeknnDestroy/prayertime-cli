## prayertime-cli times countdown

Calculate seconds and minutes remaining until a target prayer

```
prayertime-cli times countdown [flags]
```

### Examples

```
prayertime-cli times countdown --query Istanbul --target next-prayer
prayertime-cli times countdown --query Istanbul --target iftar --quiet
prayertime-cli times countdown --lat 41.01384 --lon 28.94966 --target asr --at 2026-03-07T12:00:00+03:00 --json
```

### Options

```
      --at string             Optional RFC3339 timestamp to evaluate countdown from
      --country-code string   Optional ISO country code filter
  -h, --help                  help for countdown
      --lat float             Latitude coordinate
      --lon float             Longitude coordinate
      --query string          Place name to resolve before fetching prayer times
  -q, --quiet                 Emit only remaining seconds
      --target string         Target prayer: next-prayer, imsak, fajr, sunrise, dhuhr, asr, maghrib, sunset, isha, iftar
```

### Options inherited from parent commands

```
      --json   Emit JSON output to stdout
```

### SEE ALSO

* [prayertime-cli times](prayertime-cli_times.md)	 - Fetch daily prayer times and countdowns

