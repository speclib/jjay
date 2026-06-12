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

#### Scenario: Integration target uses the formatter
- **WHEN** `make test-integration` is executed and `gotestsum` is on `PATH`
- **THEN** the integration suite runs under `gotestsum` with a verbose format
- **THEN** each scenario reports a colored PASS/FAIL line
- **THEN** a run summary reporting the test count and total duration is printed

#### Scenario: Integration target fallback
- **WHEN** `make test-integration` is executed and `gotestsum` is NOT on `PATH`
- **THEN** the integration suite still runs via `go test -tags integration ./...`
- **THEN** the run reports pass/fail results without erroring on the missing formatter

### Requirement: Test file convention
Tests SHALL be placed alongside source files using Go's `*_test.go` convention in the same package. Integration tests SHALL be placed in `test/integration/` with a `//go:build integration` tag.

#### Scenario: Initial test exists
- **WHEN** the project is scaffolded
- **THEN** at least one `*_test.go` file exists and passes

### Requirement: Full lifecycle integration test
A Go integration test (build tag `integration`) SHALL test the complete spawn → verify → cleanup → verify lifecycle using a fake agent, dedicated tmux session, and temporary jj repository.

#### Scenario: Spawn creates all resources
- **WHEN** the integration test runs `Spawn()` with a fake agent
- **THEN** a jj workspace exists with the change name
- **THEN** a tmux window exists in the test session with name `ws-<change>`
- **THEN** the workspace directory exists and contains project files
- **THEN** the fake agent's output file exists in the workspace

#### Scenario: Cleanup removes all resources
- **WHEN** the integration test runs `Cleanup()` after spawn
- **THEN** the tmux window no longer exists
- **THEN** the jj workspace no longer exists
- **THEN** the workspace directory no longer exists

### Requirement: Integration test isolation
The integration test SHALL create its own tmux session (e.g., `jjay-test-<random>`) and temporary jj repository, and SHALL clean up all resources at teardown on a normal pass/fail exit. Because teardown is bypassed on a panic or interrupt, `make test-integration` SHALL first run `clean-tests` to sweep debris left by a prior aborted run, and `make clean-tests` SHALL be available as the manual recovery path.

#### Scenario: No pollution on normal exit
- **WHEN** the integration test runs and completes (pass or fail)
- **THEN** no test tmux sessions remain
- **THEN** no test jj workspaces remain
- **THEN** no test directories remain

#### Scenario: Recovery after an aborted run
- **WHEN** a prior integration run was interrupted and left `jjay-test-*` sessions/dirs behind
- **THEN** running `make test-integration` removes them before starting (via its `clean-tests` prerequisite)
- **THEN** `make clean-tests` can also be run on its own to remove them

### Requirement: Fake agent
A `testdata/fake-agent.sh` script SHALL accept arguments, create a marker file in the current directory, and exit. This allows the integration test to verify the agent was launched in the correct directory.

#### Scenario: Fake agent runs
- **WHEN** the fake agent is launched in a workspace directory
- **THEN** it creates `agent-was-here.txt` in that directory
- **THEN** it exits with code 0

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
