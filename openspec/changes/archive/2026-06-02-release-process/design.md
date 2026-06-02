## Context

jjay is a Go CLI with cobra, built via Nix flake. Version is currently hardcoded as `"dev"` in main.go and `"0.0.1"` in flake.nix. Based on teejay's proven release pattern (goreleaser + release script + vendorHash update).

## Goals / Non-Goals

**Goals:**
- Single source of truth for version (VERSION file)
- Automated multi-platform releases via goreleaser + GitHub Actions
- Interactive release script with nix vendorHash auto-update
- Maintainer documentation

**Non-Goals:**
- Package manager distribution (homebrew, nixpkgs — future)
- Binary signing
- Automated changelog generation from commits
- Docker images

## Decisions

### VERSION file + go:embed

Plain text `VERSION` file (e.g., `0.1.0`). Go binary reads it via `//go:embed VERSION`. Flake.nix reads it via `builtins.readFile ./VERSION`. goreleaser overrides via ldflags from git tag.

_Alternative: hardcoded var + ldflags only — rejected because version drifts across flake.nix and main.go._

### goreleaser for release builds

De facto standard for Go. Handles cross-compilation, checksums, GitHub release creation. Single config file.

_Alternative: manual GitHub Actions matrix — more complex, fewer features._

### Target platforms

linux/amd64, linux/arm64, darwin/amd64, darwin/arm64. No Windows (tmux is Unix-only).

### Release script with gum

Bash script using gum for interactive version bump selection. Flow:

1. Safety checks (clean tree, on main, CHANGELOG has Unreleased)
2. Prompt for bump type (major/minor/patch)
3. Update VERSION file
4. Update CHANGELOG.md (Unreleased → version + date)
5. Compute and update vendorHash in flake.nix (skip if nix not installed)
6. Git commit, tag, push

_Alternative: Makefile target — rejected because interactive prompts are awkward in make._

### vendorHash update in release script

Temporarily set vendorHash to fake hash, run `nix build`, parse correct hash from stderr, update flake.nix. Standard nix maintainer pattern. Skip gracefully if nix not installed.

### Semantic versioning with v prefix on tags

Tags: `v0.1.0`, `v1.0.0`. VERSION file: `0.1.0` (no prefix). goreleaser expects the `v` prefix on tags.

## Risks / Trade-offs

- [flake.nix version may lag] → Release script updates it; between releases it shows last release version. Acceptable.
- [gum not installed] → Release script checks and shows install instructions.
- [nix not installed] → vendorHash update skipped with warning. Non-blocking.
- [Manual changelog] → Worth it for quality. Releases are infrequent.
