### Requirement: Nix build
The project SHALL be buildable via `nix build`, producing a single jjay binary. External tools (jj, tmux) are runtime dependencies, not build dependencies.

#### Scenario: Nix build succeeds
- **WHEN** `nix build` is executed
- **THEN** a jjay binary is produced in `result/bin/jjay`

### Requirement: Nix run
The project SHALL be runnable via `nix run`.

#### Scenario: Nix run succeeds
- **WHEN** `nix run . -- version` is executed
- **THEN** jjay prints its version string

### Requirement: Development shell
The project SHALL provide a `nix develop` shell with Go tooling (go, gopls).

#### Scenario: Dev shell has Go
- **WHEN** `nix develop -c go version` is executed
- **THEN** a Go version is printed

### Requirement: Multi-platform support
The flake SHALL support x86_64-linux, aarch64-linux, x86_64-darwin, and aarch64-darwin.

#### Scenario: Multiple systems defined
- **WHEN** `flake.nix` is inspected
- **THEN** all four system targets are defined
