## prayertime-cli

CLI-first, agent-native Islamic prayer times tool

### Synopsis

CLI-first, agent-native Islamic prayer times tool.

Contract:
  - Structured payloads go to stdout.
  - Human-readable errors and suggestions go to stderr.
  - With --output json, errors are emitted as JSON on stdout.

Global output modes:
  - --output text prints human-readable output.
  - --output json prints structured JSON.
  - --output value is reserved for commands that expose --field selectors.
  - --json is a shortcut for --output json.

Exit codes:
  - 0 success
  - 1 internal failure
  - 2 usage error
  - 3 not found or ambiguous input
  - 4 network or upstream timeout
  - 5 reserved conflict/state error

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

