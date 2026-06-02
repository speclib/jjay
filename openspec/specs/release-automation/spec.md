## ADDED Requirements

### Requirement: goreleaser builds multi-platform binaries
goreleaser SHALL build binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64. It SHALL inject the version from the git tag via ldflags. External tools (jj, tmux) are runtime dependencies not bundled.

#### Scenario: goreleaser config valid
- **WHEN** `goreleaser check` is run
- **THEN** it passes with no errors

#### Scenario: Snapshot build succeeds
- **WHEN** `goreleaser build --snapshot --clean` is run
- **THEN** binaries are produced for all four platform targets

### Requirement: GitHub Actions triggers on version tags
A GitHub Actions workflow SHALL trigger on `v*` tag pushes and run goreleaser to create a GitHub release with binaries and checksums.

#### Scenario: Workflow triggers on tag
- **WHEN** a tag matching `v*` is pushed
- **THEN** the release workflow runs goreleaser

### Requirement: Release script automates the process
A `scripts/release.sh` script SHALL provide interactive version bump selection (major/minor/patch via gum), update VERSION, update CHANGELOG.md, compute and update nix vendorHash, create git commit and tag, and push.

#### Scenario: Release script runs
- **WHEN** `scripts/release.sh` is executed on a clean main branch
- **THEN** it prompts for version bump type
- **THEN** it updates VERSION, CHANGELOG.md, and flake.nix vendorHash
- **THEN** it creates a commit, tag, and pushes

#### Scenario: Dirty working directory
- **WHEN** `scripts/release.sh` is executed with uncommitted changes
- **THEN** it exits with an error before making any changes

### Requirement: Release script updates nix vendorHash
The release script SHALL compute the correct vendorHash for `flake.nix` by attempting a nix build with a fake hash and parsing the expected hash from stderr. If nix is not installed, it SHALL warn and skip.

#### Scenario: Nix available
- **WHEN** the release script runs and nix is installed
- **THEN** flake.nix vendorHash is updated to the correct value

#### Scenario: Nix not available
- **WHEN** the release script runs and nix is not installed
- **THEN** a warning is printed and the hash update is skipped

### Requirement: Maintainer documentation
A `RELEASING.md` file SHALL document the release process: prerequisites, pre-release checklist, how to run the script, how to verify the release, and troubleshooting.

#### Scenario: Documentation exists
- **WHEN** RELEASING.md is read
- **THEN** it contains pre-release checklist, release steps, and troubleshooting

## MODIFIED Requirements

### Requirement: Nix build
The flake.nix SHALL read the version from the `VERSION` file via `builtins.readFile`. The devShell SHALL include goreleaser and gum.

#### Scenario: Nix build succeeds
- **WHEN** `nix build` is executed
- **THEN** a jjay binary is produced in `result/bin/jjay`

#### Scenario: Dev shell has release tools
- **WHEN** `nix develop` is entered
- **THEN** goreleaser and gum are available
