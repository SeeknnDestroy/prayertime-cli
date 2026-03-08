## prayertime-cli locations search

Search for a place and return candidate coordinates

### Synopsis

Search for a place before requesting prayer times.

Use this command first when a user-supplied place may be ambiguous, incomplete, or misspelled.

```
prayertime-cli locations search [flags]
```

### Examples

```
prayertime-cli locations search --query Istanbul
prayertime-cli locations search --query Springfield --country-code US --limit 3 --json
prayertime-cli locations search --query Istnbul --json
```

### Options

```
      --country-code string   Optional ISO country code filter
  -h, --help                  help for search
      --limit int             Maximum number of results to return (default 5)
      --query string          Place name to search
```

### Options inherited from parent commands

```
      --json   Emit structured JSON to stdout
```

### SEE ALSO

* [prayertime-cli locations](prayertime-cli_locations.md)	 - Search locations before requesting prayer times

