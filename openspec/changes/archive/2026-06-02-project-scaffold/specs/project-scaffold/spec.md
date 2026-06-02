## ADDED Requirements

### Requirement: Go module initialization
The project SHALL have a `go.mod` file at the project root with Go 1.24+ as minimum version.

#### Scenario: Module exists
- **WHEN** the project root is inspected
- **THEN** `go.mod` exists with a valid module path and Go version >= 1.24

### Requirement: CLI entry point with cobra
The project SHALL have a cobra-based CLI entry point at `cmd/jjay/main.go`. The root command SHALL display help text describing jjay.

#### Scenario: Binary runs
- **WHEN** `go run ./cmd/jjay` is executed
- **THEN** jjay displays help text and exits cleanly

### Requirement: Version command
The CLI SHALL include a `version` subcommand that prints the current version string.

#### Scenario: Version output
- **WHEN** `jjay version` is executed
- **THEN** a version string is printed to stdout

### Requirement: Directory structure
The project SHALL follow Go conventions with `cmd/jjay/` for the CLI entry point and `internal/` for private packages.

#### Scenario: Structure exists
- **WHEN** the project layout is inspected
- **THEN** `cmd/jjay/` and `internal/` directories exist
