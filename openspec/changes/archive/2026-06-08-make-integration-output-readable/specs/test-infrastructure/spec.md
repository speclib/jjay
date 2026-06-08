## MODIFIED Requirements

### Requirement: Makefile with dev targets
The project SHALL have a Makefile with `build`, `test`, `lint`, `coverage`, `badge`, `test-spawn`, and `test-integration` targets. The `test-integration` target SHALL format its output with `gotestsum` when available, using a verbose format so subprocess detail remains visible, and SHALL fall back to plain `go test` when `gotestsum` is not installed.

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

#### Scenario: Integration target uses the formatter
- **WHEN** `make test-integration` is executed and `gotestsum` is on `PATH`
- **THEN** the integration suite runs under `gotestsum` with a verbose format
- **THEN** each scenario reports a colored PASS/FAIL line
- **THEN** a run summary reporting the test count and total duration is printed

#### Scenario: Integration target fallback
- **WHEN** `make test-integration` is executed and `gotestsum` is NOT on `PATH`
- **THEN** the integration suite still runs via `go test -tags integration ./...`
- **THEN** the run reports pass/fail results without erroring on the missing formatter

## ADDED Requirements

### Requirement: Readable integration test output
The integration test suite SHALL produce output in which each scenario is visually distinguished by a heading line, subprocess detail (jj, OpenSpec, tmux) is nested beneath the scenario it belongs to rather than streamed as free-floating text, and a run summary with pass/fail counts and total duration is emitted. Per-scenario results SHALL be colorized when the formatter and terminal support color.

#### Scenario: Subprocess output is attributed to its scenario
- **WHEN** an integration test helper shells out to `jj`, `openspec`, or `tmux`
- **THEN** that subprocess's combined stdout/stderr is emitted through the test logger
- **THEN** the captured output appears under the heading of the test (or subtest) that invoked it, not as unattributed lines in the terminal

#### Scenario: Run summary is present
- **WHEN** the integration suite finishes under the formatter
- **THEN** a summary line reporting the number of tests run and the total elapsed time is printed at the end of the run

#### Scenario: Failures remain diagnosable
- **WHEN** an integration scenario fails
- **THEN** the failing scenario is clearly marked (colored FAIL where supported)
- **THEN** the subprocess detail captured for that scenario is visible in the output
