## ADDED Requirements

### Requirement: Cleanup kills tmux window
The `jjay cleanup <change>` command SHALL kill the tmux window named `ws-<change>` if it exists. If the window does not exist, it SHALL skip this step without error.

#### Scenario: Window exists
- **WHEN** `jjay cleanup feat-payments` is executed and tmux window `ws-feat-payments` exists
- **THEN** the window is killed

#### Scenario: Window already gone
- **WHEN** `jjay cleanup feat-payments` is executed and tmux window `ws-feat-payments` does not exist
- **THEN** the step is skipped
- **THEN** cleanup continues without error

### Requirement: Cleanup forgets jj workspace
The command SHALL forget the jj workspace named `<change>` if it exists. If the workspace does not exist, it SHALL skip this step without error.

#### Scenario: Workspace exists
- **WHEN** `jjay cleanup feat-payments` is executed and jj workspace `feat-payments` exists
- **THEN** `jj workspace forget feat-payments` is run

#### Scenario: Workspace already gone
- **WHEN** `jjay cleanup feat-payments` is executed and jj workspace `feat-payments` does not exist
- **THEN** the step is skipped
- **THEN** cleanup continues without error

### Requirement: Cleanup removes workspace directory
The command SHALL remove the workspace directory at `../<project-name>-workspaces/<change>` if it exists. If the directory does not exist, it SHALL skip this step without error.

#### Scenario: Directory exists
- **WHEN** `jjay cleanup feat-payments` is executed in project `jjay` and `../jjay-workspaces/feat-payments` exists
- **THEN** the directory and all contents are removed

#### Scenario: Directory already gone
- **WHEN** `jjay cleanup feat-payments` is executed and the workspace directory does not exist
- **THEN** the step is skipped
- **THEN** cleanup continues without error

### Requirement: Cleanup reports what it did
The command SHALL print a summary of which steps were performed and which were skipped.

#### Scenario: Full cleanup
- **WHEN** all three resources exist (window, workspace, directory)
- **THEN** output indicates all three were cleaned up

#### Scenario: Partial cleanup
- **WHEN** only the directory exists (window and workspace already gone)
- **THEN** output indicates the directory was removed and the other steps were skipped

### Requirement: Cleanup execution order
The command SHALL execute in this order: kill tmux window first (stops the running agent), then forget jj workspace, then remove directory.

#### Scenario: Order preserved
- **WHEN** `jjay cleanup feat-payments` is executed
- **THEN** tmux kill-window runs before jj workspace forget
- **THEN** jj workspace forget runs before directory removal

### Requirement: Cleanup requires change name argument
The command SHALL require exactly one argument: the change name.

#### Scenario: No argument
- **WHEN** `jjay cleanup` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code

### Requirement: Cleanup has unit-test coverage
The `internal/cleanup` package SHALL have unit tests covering its testable behavior: workspace-directory removal (present and absent), the tmux target helper, and the tolerance branches that skip missing tmux/jj resources without error. Tests SHALL NOT require a live tmux server or a jj repository.

#### Scenario: Removes an existing workspace directory
- **WHEN** a workspace directory exists and `removeDirectory` runs for its change (with the matching workspace root)
- **THEN** the directory no longer exists afterward

#### Scenario: Tolerates a missing workspace directory
- **WHEN** the workspace directory does not exist and `removeDirectory` runs
- **THEN** it returns without error and reports the directory was skipped

#### Scenario: tmux target formatting
- **WHEN** the tmux target is built with no session
- **THEN** it equals the window name
- **WHEN** the tmux target is built with a session
- **THEN** it equals `<session>:<window>`

#### Scenario: Cleanup is tolerant of all-missing resources
- **WHEN** `Cleanup` runs for a change whose tmux window, jj workspace, and directory are all absent
- **THEN** it returns nil (no error)
