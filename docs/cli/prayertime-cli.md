## prayertime-cli

CLI-first, agent-native Islamic prayer times tool

### Synopsis

CLI-first, agent-native Islamic prayer times tool.

Start here:
  - Use 'times get' for one day's prayer schedule.
  - Use 'times countdown' to ask how long until the next prayer or a named prayer.
  - Use 'locations search' first if a place name may be ambiguous.

Common tasks:
  - prayertime-cli times get --query Istanbul --json
  - prayertime-cli times countdown --query Istanbul --target next-prayer --json
  - prayertime-cli locations search --query Springfield --country-code US --json

Output model:
  - Structured payloads go to stdout.
  - Human-readable errors and suggestions go to stderr.
  - --json is a shortcut for structured JSON output.
  - --output text|json|value is the generalized output switch.

### Options

```
  -h, --help            help for prayertime-cli
      --json            Shortcut for --output json
      --output string   Output mode: text, json, or value (default "text")
```

### SEE ALSO

* [prayertime-cli locations](prayertime-cli_locations.md)	 - Search locations before requesting prayer times
* [prayertime-cli times](prayertime-cli_times.md)	 - Fetch daily prayer times and countdowns
* [prayertime-cli version](prayertime-cli_version.md)	 - Print the CLI version

