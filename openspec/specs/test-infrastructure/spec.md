### Requirement: Makefile with dev targets
The project SHALL have a Makefile with `build`, `test`, `lint`, `coverage`, `badge`, and `clean-tests` targets.

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

#### Scenario: Clean-tests target
- **WHEN** `make clean-tests` is executed
- **THEN** every tmux session whose name matches `jjay-test-*` is killed
- **THEN** temp directories `/tmp/jjay-test-*` and `/tmp/jjay-merge-test-*` are removed
- **THEN** real spawn sessions (named `jjay-><dirname>`) are NOT affected

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
