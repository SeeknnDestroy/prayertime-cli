## prayertime-cli times

Fetch daily prayer times and countdowns

### Synopsis

Fetch daily prayer times and countdowns for a resolved location.

Contract:
  - Provide exactly one location input strategy: --query <place> or --lat with --lon.
  - --output text prints human-readable output to stdout.
  - --output json prints structured JSON to stdout.
  - --output value prints only the selected --field value.
  - Human-readable errors go to stderr. With --output json, errors are JSON on stdout.

### Options

```
  -h, --help   help for times
```

### Options inherited from parent commands

```
      --json            Shortcut for --output json
      --output string   Output mode: text, json, or value (default "text")
```

### SEE ALSO

* [prayertime-cli](prayertime-cli.md)	 - CLI-first, agent-native Islamic prayer times tool
* [prayertime-cli times countdown](prayertime-cli_times_countdown.md)	 - Calculate seconds and minutes remaining until a target prayer
* [prayertime-cli times get](prayertime-cli_times_get.md)	 - Fetch prayer times for a location and date

