## prayertime-cli locations search

Search for a place and return candidate coordinates

```
prayertime-cli locations search [flags]
```

### Examples

```
prayertime-cli locations search --query Istanbul
prayertime-cli locations search --query London --country-code GB --json
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
      --json   Emit JSON output to stdout
```

### SEE ALSO

* [prayertime-cli locations](prayertime-cli_locations.md)	 - Search locations before requesting prayer times

