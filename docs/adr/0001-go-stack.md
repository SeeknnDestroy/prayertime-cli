# ADR 0001: Use Go For The CLI

## Status

Accepted

## Context

The project needs a fast, cross-platform, single-binary CLI with strong standard-library networking support and low runtime overhead.

## Decision

Use Go 1.26 as the implementation language and Cobra for the command tree.

## Consequences

- Distribution is straightforward through GitHub Releases, Homebrew, and Scoop.
- The standard library covers most of the runtime surface for MVP 1.
- Contributors need a Go toolchain instead of Python or Rust.

