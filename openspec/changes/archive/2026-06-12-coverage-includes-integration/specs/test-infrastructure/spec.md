## MODIFIED Requirements

### Requirement: Makefile with dev targets
The project SHALL have a Makefile with `build`, `test`, `lint`, `coverage`, `coverage-unit`, `badge`, `test-spawn`, `test-integration`, and `clean-tests` targets. The `test-integration` target SHALL format its output with `gotestsum` when available, using a verbose format so subprocess detail remains visible, and SHALL fall back to plain `go test` when `gotestsum` is not installed.

The `coverage` target SHALL measure coverage across the whole repository including the integration tests: it SHALL run `go test` with `-tags integration` and `-coverpkg=./...` so that (a) integration-tested code is exercised and (b) coverage of every package is attributed regardless of which test package drove it (tests in `test/integration` exercise `internal/spawn`, `internal/cleanup`, etc. in-process but from a different package, and are not attributed without `-coverpkg`). Because it runs the integration suite, `coverage` SHALL sweep test debris first (depend on `clean-tests`) and SHALL require `tmux` and `jj` on `PATH`. A `coverage-unit` target SHALL provide the same whole-repo `-coverpkg` measurement WITHOUT `-tags integration`, for environments lacking tmux/jj (bare CI).

The README coverage badge SHALL be updated only by `make badge` (which depends on `coverage`); `make coverage` itself only computes and prints the percentage.

#### Scenario: Build target
- **WHEN** `make build` is executed
- **THEN** a `jjay` binary is produced

#### Scenario: Test target
- **WHEN** `make test` is executed
- **THEN** `go test ./...` runs and reports results

#### Scenario: Lint target
- **WHEN** `make lint` is executed
- **THEN** `go vet ./...` runs and reports results

#### Scenario: Coverage target includes integration with whole-repo attribution
- **WHEN** `make coverage` is executed (with tmux and jj on `PATH`)
- **THEN** test debris is swept first (via `clean-tests`)
- **THEN** the integration suite runs with `-tags integration -coverpkg=./...`
- **THEN** integration-exercised packages (e.g. `internal/merge`, `internal/spawn`) report their true coverage, not ~5%
- **THEN** an HTML coverage report is generated and the total percentage is printed to stdout

#### Scenario: Coverage-unit fallback without tmux/jj
- **WHEN** `make coverage-unit` is executed
- **THEN** coverage runs with `-coverpkg=./...` but WITHOUT `-tags integration`
- **THEN** it completes without requiring tmux or jj (integration-only coverage is simply absent)

#### Scenario: Badge target updates the README
- **WHEN** `make badge` is executed
- **THEN** `coverage` runs first and the README coverage badge is updated with the current percentage
- **THEN** running `make coverage` alone does NOT modify the README

#### Scenario: Clean-tests target
- **WHEN** `make clean-tests` is executed
- **THEN** every tmux session whose name matches `jjay-test-*` is killed
- **THEN** temp directories `/tmp/jjay-test-*` and `/tmp/jjay-merge-test-*` are removed
- **THEN** real spawn sessions (named `jjay-><dirname>`) are NOT affected
