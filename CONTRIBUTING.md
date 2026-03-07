# Contributing

## Workflow

1. Create a feature branch.
2. Keep commits atomic and use Conventional Commit messages.
3. Run `go test ./...` before opening or updating a pull request.
4. Open a pull request instead of pushing directly to `main`.

## Standards

- Preserve the documented CLI contract and exit codes.
- Keep JSON output stable once released.
- Prefer focused changes with tests over broad refactors.

