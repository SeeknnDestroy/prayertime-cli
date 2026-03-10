# Release Workflow

## Required GitHub Variables

- `HOMEBREW_TAP_OWNER`
- `HOMEBREW_TAP_NAME`
- `SCOOP_BUCKET_OWNER`
- `SCOOP_BUCKET_NAME`

Current companion repositories:

- `SeeknnDestroy/homebrew-tap`
- `SeeknnDestroy/scoop-bucket`

Homebrew publishing uses a cask in the tap's `Casks/` directory.

## Required GitHub Secret

- `GH_PAT`

`GH_PAT` needs `repo` scope because GoReleaser cannot push Homebrew and Scoop updates to other repositories with the default `GITHUB_TOKEN`.

## Release Commands

```bash
make verify
make release-check
make release-snapshot
git tag v0.1.0
git push origin v0.1.0
```
