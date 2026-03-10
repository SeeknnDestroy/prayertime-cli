# Contributing

## Workflow

1. Create a feature branch.
2. Keep commits atomic and use Conventional Commit messages.
3. Use the same Conventional Commit style for pull request titles: `type(scope): subject` or `type: subject`.
4. Run `make verify` before opening or updating a pull request.
5. Open a pull request instead of pushing directly to `main`.

## Standards

- Preserve the documented CLI contract and exit codes.
- Keep JSON output stable once released.
- Keep `--json` and `--quiet` discoverable for agent-facing workflows.
- Prefer focused changes with tests over broad refactors.
