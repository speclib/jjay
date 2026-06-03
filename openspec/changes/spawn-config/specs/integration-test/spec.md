## ADDED Requirements

### Requirement: Full lifecycle integration test
A Go integration test (build tag `integration`) SHALL test the complete spawn → verify → cleanup → verify lifecycle using a fake agent, dedicated tmux session, and temporary jj repository.

#### Scenario: Spawn creates all resources
- **WHEN** the integration test runs `Spawn()` with a fake agent
- **THEN** a jj workspace exists with the change name
- **THEN** a tmux window exists in the test session with name `ws-<change>`
- **THEN** the workspace directory exists and contains project files
- **THEN** the fake agent's output file exists in the workspace

#### Scenario: Panes are in correct working directory
- **WHEN** the integration test inspects the tmux panes after spawn
- **THEN** both panes have their working directory set to the workspace directory
- **THEN** this is verified via `tmux display-message -p -t <pane> '#{pane_current_path}'`

#### Scenario: Cleanup removes all resources
- **WHEN** the integration test runs `Cleanup()` after spawn
- **THEN** the tmux window no longer exists
- **THEN** the jj workspace no longer exists
- **THEN** the workspace directory no longer exists

### Requirement: Integration test isolation
The integration test SHALL create its own tmux session (e.g., `jjay-test-<random>`) and temporary jj repository. It SHALL clean up all resources at teardown regardless of pass/fail.

#### Scenario: No pollution
- **WHEN** the integration test runs and completes (pass or fail)
- **THEN** no test tmux sessions remain
- **THEN** no test jj workspaces remain
- **THEN** no test directories remain

### Requirement: Fake agent
A `testdata/fake-agent.sh` script SHALL accept arguments, create a marker file in the current directory, and exit. This allows the integration test to verify the agent was launched in the correct directory.

#### Scenario: Fake agent runs
- **WHEN** the fake agent is launched in a workspace directory
- **THEN** it creates `agent-was-here.txt` in that directory
- **THEN** it exits with code 0
