## ADDED Requirements

### Requirement: Session-open reopens detached spawns
After creating and switching to the tmux session, `jjay session-open <path>` SHALL recreate a `ws-<change>` tmux window for every spawned jj workspace in the repository that does not already have a matching window in the new session. The recreated window SHALL be created in the workspace directory and SHALL relaunch the agent the same way `jjay spawn` does.

#### Scenario: Detached spawns are reopened
- **WHEN** the repository has jj workspaces `add-foo` and `fix-bar` and `jjay session-open <path>` creates a fresh session with no windows
- **THEN** a `ws-add-foo` window and a `ws-fix-bar` window are created in the session
- **THEN** each window's working directory is the corresponding workspace directory
- **THEN** the agent is launched in each window

#### Scenario: No spawns to reopen
- **WHEN** the repository has no spawned jj workspaces and `jjay session-open <path>` is executed
- **THEN** the session is created with no `ws-*` windows
- **THEN** jjay exits zero

#### Scenario: Existing window is not duplicated
- **WHEN** a `ws-add-foo` window already exists in the session and reopen runs
- **THEN** no second `ws-add-foo` window is created

### Requirement: Reopen failures do not abort session-open
If recreating a window for one spawn fails, `jjay session-open` SHALL continue reopening the remaining spawns and SHALL still complete session creation, reporting which spawns failed to reopen.

#### Scenario: One reopen fails
- **WHEN** reopening `add-foo` fails but `fix-bar` succeeds
- **THEN** the `ws-fix-bar` window exists
- **THEN** the session-open command completes (non-fatal) and reports that `add-foo` could not be reopened
