## ADDED Requirements

### Requirement: VERSION file as single source of truth
A `VERSION` file at the project root SHALL contain the version number (e.g., `0.1.0`). No `v` prefix.

#### Scenario: VERSION file exists
- **WHEN** the project root is inspected
- **THEN** a `VERSION` file exists with a valid semver string

### Requirement: Binary embeds version from VERSION file
The `jjay version` command SHALL print the version from the embedded VERSION file. At build time, `go:embed` reads the file contents.

#### Scenario: Dev build
- **WHEN** `go run ./cmd/jjay version` is executed
- **THEN** the version from the VERSION file is printed

#### Scenario: ldflags override
- **WHEN** the binary is built with `-ldflags "-X main.version=v1.2.3"`
- **THEN** `jjay version` prints `v1.2.3` (ldflags takes precedence for goreleaser)
