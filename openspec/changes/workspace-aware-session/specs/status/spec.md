## ADDED Requirements

### Requirement: Status lists jj workspaces
The `jjay status` command SHALL list every jj workspace in the current repository, derived live from `jj workspace list`. It SHALL persist no state and read no state file. The default jjay workspace (the main working copy) MAY be shown but SHALL be distinguishable from spawned workspaces.

#### Scenario: Spawned workspaces are listed
- **WHEN** two spawns exist (jj workspaces `add-foo` and `fix-bar`) and `jjay status` is executed
- **THEN** the output includes a row for `add-foo` and a row for `fix-bar`
- **THEN** each row shows the change name and its workspace directory

#### Scenario: No spawns
- **WHEN** no spawned jj workspaces exist and `jjay status` is executed
- **THEN** jjay exits zero
- **THEN** the output indicates there are no running spawns

### Requirement: Status reports attached vs detached
For each spawned workspace, the command SHALL report whether a matching tmux window named `ws-<change>` exists in the current tmux session. A workspace WITH a matching window SHALL be reported as **attached**; a workspace WITHOUT one SHALL be reported as **detached**. A detached workspace is still open — its absence of a window SHALL NOT be reported as closed or errored.

#### Scenario: Attached spawn
- **WHEN** workspace `add-foo` exists and a `ws-add-foo` window exists in the current session, and `jjay status` is executed
- **THEN** the `add-foo` row is reported as attached

#### Scenario: Detached spawn
- **WHEN** workspace `add-foo` exists but no `ws-add-foo` window exists in the current session, and `jjay status` is executed
- **THEN** the `add-foo` row is reported as detached
- **THEN** jjay exits zero

### Requirement: Status detection is scoped to the current session
Window detection SHALL be scoped to the current tmux session. Cross-session discovery is out of scope.

#### Scenario: Window in another session is not counted as attached
- **WHEN** a `ws-add-foo` window exists only in a different tmux session, and `jjay status` is executed from the current session
- **THEN** the `add-foo` row is reported as detached

### Requirement: Status tolerates missing tmux
The command SHALL NOT fail when no tmux server or session is available; in that case all spawned workspaces SHALL be reported as detached.

#### Scenario: No tmux server
- **WHEN** no tmux server is running and `jjay status` is executed
- **THEN** jjay exits zero
- **THEN** every spawned workspace is reported as detached

### Requirement: Status takes no arguments
The command SHALL accept no positional arguments.

#### Scenario: Unexpected argument
- **WHEN** `jjay status extra-arg` is executed
- **THEN** cobra prints usage help and exits with non-zero exit code
