## ADDED Requirements

### Requirement: Session-open creates tmux session
The `jjay session-open <path>` command SHALL create a new tmux session named `jjay-><dirname>` where `<dirname>` is the basename of the given path. The session's working directory SHALL be set to the given path.

#### Scenario: Session created
- **WHEN** `jjay session-open ~/projects/myapp` is executed
- **THEN** a tmux session named `jjay->myapp` is created
- **THEN** the session's working directory is `~/projects/myapp`

### Requirement: Session-open switches to the session
After creating the session, the command SHALL switch the current tmux client to the new session.

#### Scenario: Client switches
- **WHEN** `jjay session-open ~/projects/myapp` completes
- **THEN** the user is in the `jjay->myapp` tmux session

### Requirement: Session-open requires jj repo
The command SHALL verify the given path contains a jj repository (`.jj/` directory exists).

#### Scenario: Valid jj repo
- **WHEN** `jjay session-open ~/projects/myapp` is executed and `~/projects/myapp/.jj/` exists
- **THEN** session creation proceeds

#### Scenario: Not a jj repo
- **WHEN** `jjay session-open /tmp/notarepo` is executed and no `.jj/` directory exists
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the path is not a jj repository

### Requirement: Session-open rejects duplicate sessions
The command SHALL check if a session with that name already exists.

#### Scenario: Session already exists
- **WHEN** `jjay session-open ~/projects/myapp` is executed and `jjay->myapp` session exists
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the session already exists

### Requirement: Session-open requires path argument
The command SHALL require exactly one argument: the path to the repo.

#### Scenario: No argument
- **WHEN** `jjay session-open` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code

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
