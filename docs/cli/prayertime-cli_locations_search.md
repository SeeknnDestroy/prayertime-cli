## prayertime-cli locations search

Search for a place and return candidate coordinates

### Synopsis

Resolve a place query into candidate locations and canonical coordinates.

Input rules:
  - --query is required.
  - --country-code is optional and should be used to narrow ambiguous place names.
  - --limit controls the maximum number of returned candidates.

Output:
  - --output text prints numbered candidates with coordinates and timezone.
  - --output json prints query, count, and structured candidates including display_name.
  - --output value is not supported for location search.

```
prayertime-cli locations search [flags]
```

### Examples

```
prayertime-cli locations search --query Istanbul
prayertime-cli locations search --query London --country-code GB --output json
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
      --json            Shortcut for --output json
      --output string   Output mode: text, json, or value (default "text")
```

### SEE ALSO

* [prayertime-cli locations](prayertime-cli_locations.md)	 - Search locations before requesting prayer times

