## ADDED Requirements

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
