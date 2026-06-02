# ADR-004: VERSION file as single source of truth

**Status**: Proposed

## Context

Version is currently hardcoded in two places: `"dev"` in `cmd/jjay/main.go` and `"0.0.1"` in `flake.nix`. This will drift. goreleaser adds a third source (git tag). Need one canonical location.

## Options Considered

- **VERSION file + go:embed + builtins.readFile** — single file, read at build time by both Go and Nix. goreleaser overrides via ldflags from git tag.
- **ldflags only** — goreleaser sets version at build time. But flake.nix can't read ldflags, so it still needs its own source.
- **Git tag as source** — clean but requires tag to exist before build. Doesn't work for dev builds.

## Decision

Plain text `VERSION` file at project root. Go reads it via `//go:embed`, Nix reads it via `builtins.readFile ./VERSION`. goreleaser overrides via ldflags from the git tag (which should match). Release script keeps VERSION file in sync.

## Consequences

- **Positive**: One file to update, two consumers read it automatically
- **Positive**: Dev builds show the current version, not "dev"
- **Negative**: VERSION file must be kept in sync with git tags — release script handles this
