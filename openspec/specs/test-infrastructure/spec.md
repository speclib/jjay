### Requirement: Makefile with dev targets
The project SHALL have a Makefile with `build`, `test`, and `lint` targets.

#### Scenario: Build target
- **WHEN** `make build` is executed
- **THEN** a `jjay` binary is produced

#### Scenario: Test target
- **WHEN** `make test` is executed
- **THEN** `go test ./...` runs and reports results

#### Scenario: Lint target
- **WHEN** `make lint` is executed
- **THEN** `go vet ./...` runs and reports results

### Requirement: Test file convention
Tests SHALL be placed alongside source files using Go's `*_test.go` convention in the same package.

#### Scenario: Initial test exists
- **WHEN** the project is scaffolded
- **THEN** at least one `*_test.go` file exists and passes
