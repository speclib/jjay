## MODIFIED Requirements

### Requirement: Makefile with dev targets
The project SHALL have a Makefile with `build`, `test`, `lint`, `coverage`, and `badge` targets.

#### Scenario: Build target
- **WHEN** `make build` is executed
- **THEN** a `jjay` binary is produced

#### Scenario: Test target
- **WHEN** `make test` is executed
- **THEN** `go test ./...` runs and reports results

#### Scenario: Lint target
- **WHEN** `make lint` is executed
- **THEN** `go vet ./...` runs and reports results

#### Scenario: Coverage target
- **WHEN** `make coverage` is executed
- **THEN** unit tests run with coverage profiling
- **THEN** an HTML coverage report is generated
- **THEN** the total coverage percentage is printed to stdout

#### Scenario: Badge target
- **WHEN** `make badge` is executed
- **THEN** the README coverage badge is updated with the current percentage
