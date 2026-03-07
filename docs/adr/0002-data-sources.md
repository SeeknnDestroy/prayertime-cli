# ADR 0002: Use Open-Meteo Plus AlAdhan Method 13

## Status

Accepted

## Context

The CLI needs globally usable location search and prayer times aligned with Turkish Diyanet expectations.

## Decision

- Resolve place names with Open-Meteo geocoding.
- Fetch prayer times with AlAdhan using `method=13`.
- Treat coordinates as canonical and avoid relying on AlAdhan city metadata.

## Consequences

- The CLI remains stateless and API-key free in MVP 1.
- Upstream schema drift must be monitored with tests and scheduled contract checks.
- Ambiguous locations need explicit user-visible recovery paths.

