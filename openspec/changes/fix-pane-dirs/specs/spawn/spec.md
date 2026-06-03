## MODIFIED Requirements

### Requirement: Spawn creates tmux window
The command SHALL create a new tmux window in the target session named `ws-<change>` with its starting directory set to the workspace directory via `tmux new-window -c <wsDir>`.

#### Scenario: Window starts in workspace dir
- **WHEN** `jjay spawn feat-payments` is executed
- **THEN** the tmux window's starting directory is the workspace directory
- **THEN** verified via `tmux display-message -p -t <window> '#{pane_current_path}'`

### Requirement: Spawn creates two-pane layout
The tmux window SHALL be split into two panes via `tmux split-window -h -c <wsDir>`. Both panes SHALL have their working directory set to the workspace directory at creation time. The left pane SHALL run the agent command. The right pane SHALL be a plain shell. No `send-keys cd` SHALL be used for setting working directories.

#### Scenario: Both panes in workspace dir
- **WHEN** `jjay spawn feat-payments` is executed successfully
- **THEN** both pane 0 and pane 1 have working directory set to the workspace directory
- **THEN** verified via `tmux display-message -p -t <pane> '#{pane_current_path}'`

## ADDED Requirements

### Requirement: Integration test for spawn/cleanup lifecycle
A Go integration test (build tag `integration`) SHALL test the full lifecycle: spawn with fake agent → verify resources and pane directories → cleanup → verify all resources removed.

#### Scenario: Panes are in correct working directory
- **WHEN** the integration test inspects panes after spawn
- **THEN** both panes report the workspace directory as their current path

#### Scenario: Spawn creates all resources
- **WHEN** the integration test runs spawn with a fake agent
- **THEN** a jj workspace exists, tmux window exists, workspace directory exists, agent marker file exists

#### Scenario: Cleanup removes all resources
- **WHEN** the integration test runs cleanup after spawn
- **THEN** tmux window, jj workspace, and workspace directory are all gone
