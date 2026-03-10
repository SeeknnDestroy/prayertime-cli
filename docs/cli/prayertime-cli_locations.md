## prayertime-cli locations

Search locations before requesting prayer times

### Synopsis

Search and inspect candidate places before requesting prayer times.

Contract:
  - --query is required.
  - Use --country-code to narrow ambiguous place names.
  - --output text prints numbered candidates for humans.
  - --output json prints structured candidates with display_name, coordinates, and timezone.
  - Success payloads go to stdout. Errors go to stderr unless --output json is used.

### Options

```
  -h, --help   help for locations
```

### Options inherited from parent commands

```
      --json            Shortcut for --output json
      --output string   Output mode: text, json, or value (default "text")
```

### SEE ALSO

* [prayertime-cli](prayertime-cli.md)	 - CLI-first, agent-native Islamic prayer times tool
* [prayertime-cli locations search](prayertime-cli_locations_search.md)	 - Search for a place and return candidate coordinates

