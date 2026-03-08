## prayertime-cli

CLI-first, agent-native Islamic prayer times tool

### Synopsis

Stateless prayer-time CLI for agents, shell scripts, and direct terminal use.

MVP 1 has no persisted default location. Every times query must provide --query PLACE or both --lat and --lon.

### Examples

```
prayertime-cli locations search --query "Springfield" --country-code US --limit 3 --json
prayertime-cli times get --query Istanbul --json
prayertime-cli times countdown --query Istanbul --target next-prayer --json
```

### Options

```
  -h, --help   help for prayertime-cli
      --json   Emit structured JSON to stdout
```

### SEE ALSO

* [prayertime-cli locations](prayertime-cli_locations.md)	 - Search locations before requesting prayer times
* [prayertime-cli times](prayertime-cli_times.md)	 - Fetch daily prayer schedules and countdowns
* [prayertime-cli version](prayertime-cli_version.md)	 - Print the CLI version

