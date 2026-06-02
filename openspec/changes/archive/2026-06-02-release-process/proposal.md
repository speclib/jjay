# Proposal: Release process

**Change**: release-process
**Status**: proposed
**Bean**: [jjay-fder — release process](../../../.beans/jjay-fder--release-process-goreleaser-github-actions-versioni.md)

## Why

jjay needs a repeatable release process before the first release. Currently version is hardcoded as `"dev"` in `main.go` and `"0.0.1"` in `flake.nix`. No automation exists for building release binaries or creating GitHub releases.

## What Changes

- Add `VERSION` file as single source of truth for version
- Wire version embedding via `go:embed` in `cmd/jjay/main.go`
- Wire `flake.nix` to read version from `VERSION` file
- Add goreleaser config for multi-platform binary builds
- Add GitHub Actions workflow triggered by `v*` tags
- Add `scripts/release.sh` for interactive release automation (uses gum)
- Add nix vendorHash auto-update step in release script
- Add `RELEASING.md` maintainer docs

## Capabilities

### New Capabilities

- `release-automation`: goreleaser config, GitHub Actions workflow, release script with vendorHash update
- `version-embedding`: Single `VERSION` file, go:embed in binary, builtins.readFile in flake.nix

### Modified Capabilities

- `nix-flake`: Read version from VERSION file, add goreleaser to devShell

## Impact

- **New files**: `VERSION`, `.goreleaser.yaml`, `.github/workflows/release.yml`, `scripts/release.sh`, `RELEASING.md`
- **Modified files**: `cmd/jjay/main.go` (go:embed VERSION), `flake.nix` (readFile version, goreleaser + gum in devShell)
- **Build process**: Releases triggered by pushing `v*` tags
- **Dev dependencies**: goreleaser, gum (in devShell only)
