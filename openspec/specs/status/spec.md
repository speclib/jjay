## ADDED Requirements

### Requirement: Status lists jj workspaces
The `jjay status` command SHALL list every jj workspace in the current repository, derived live from `jj workspace list`. It SHALL persist no state and read no state file. The default jjay workspace (the main working copy) MAY be shown but SHALL be distinguishable from spawned workspaces.

#### Scenario: Spawned workspaces are listed
- **WHEN** two spawns exist (jj workspaces `add-foo` and `fix-bar`) and `jjay status` is executed
- **THEN** the output includes a row for `add-foo` and a row for `fix-bar`
- **THEN** each row shows the change name and its workspace directory

### Requirement: Workspace directories are shown relative to the main repo root
Each row's workspace directory SHALL be displayed as a path relative to the main jj working copy (the `default` workspace) root, not as an absolute path. The main repo root SHALL be resolved even when `jjay status` is run from inside a child workspace, by following the workspace's `.jj/repo` pointer to the main copy. If a relative path cannot be computed, the absolute path MAY be shown as a fallback.

#### Scenario: Path is relative to main root
- **WHEN** a spawn's workspace directory is `<main>/../<project>-workspaces/add-foo` and `jjay status` is executed from the main repo root
- **THEN** the `add-foo` row shows the directory relative to the main repo root (e.g. `../<project>-workspaces/add-foo`), not an absolute path

#### Scenario: Run from inside a child workspace
- **WHEN** `jjay status` is executed from within a spawned child workspace (whose `.jj/repo` is a pointer to the main copy)
- **THEN** the main repo root is resolved via the pointer
- **THEN** every workspace directory is shown relative to that main repo root

### Requirement: Status reports task progress per spawn
For each spawned workspace, the command SHALL report openspec task progress read from that workspace's `openspec/changes/<change>/tasks.md`, counting completed (`- [x]`) versus total checkboxes, rendered as `done/total (percent%)`. A workspace with no readable `tasks.md` SHALL render a placeholder (e.g. `-`) and SHALL NOT cause an error.

#### Scenario: Task counts are shown
- **WHEN** workspace `add-foo` has a `tasks.md` with 12 of 18 checkboxes marked done and `jjay status` is executed
- **THEN** the `add-foo` row shows `12/18 (66%)`

#### Scenario: Missing tasks.md
- **WHEN** workspace `add-foo` has no `openspec/changes/add-foo/tasks.md` and `jjay status` is executed
- **THEN** the `add-foo` row shows a placeholder for tasks
- **THEN** jjay exits zero

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

### Requirement: Status reports whether a spawn is merged
`jjay status` SHALL show, for each spawned workspace, whether its work has already landed on the `main` bookmark, in a `MERGED` column (yes/no). This SHALL be derived live from jj with no state file. A spawn is reported merged when its workspace has no commits that `main` lacks.

#### Scenario: Unmerged spawn
- **WHEN** a spawn's workspace has commits not present on `main` and `jjay status` runs
- **THEN** its `MERGED` value is `no`

#### Scenario: Merged spawn still on disk
- **WHEN** a spawn's work has been merged into `main` but its workspace still exists (not yet cleaned up) and `jjay status` runs
- **THEN** its `MERGED` value is `yes`

#### Scenario: Merged is independent of archived
- **WHEN** a spawn is merged into `main` but its openspec change is not yet archived
- **THEN** `MERGED` is `yes` and `ARCHIVED` is `no`

### Requirement: Status tmux-state column is named TMUX
The column reporting tmux window attached/detached state SHALL be headed `TMUX` (previously `STATUS`). The values remain `attached` / `detached`.

#### Scenario: Column header
- **WHEN** `jjay status` prints its table with at least one spawn
- **THEN** the header includes a `TMUX` column and does not use `STATUS` for the tmux state

### Requirement: Status column set
The `jjay status` table SHALL include the columns `CHANGE`, `WORKSPACE`, `TASKS`, `TMUX`, `MERGED`, and `ARCHIVED`.

#### Scenario: All columns present
- **WHEN** `jjay status` prints its table with at least one spawn
- **THEN** the header row contains CHANGE, WORKSPACE, TASKS, TMUX, MERGED, and ARCHIVED

### Requirement: Status is covered by the integration test
The lifecycle integration test (build tag `integration`) SHALL exercise `jjay status` against a real spawned workspace, asserting the column set and that a freshly-spawned, unmerged workspace reports `MERGED=no`.

#### Scenario: Status reported for a live spawn
- **WHEN** the lifecycle test has spawned a workspace and runs `status` before cleanup
- **THEN** the spawn appears in the status output
- **THEN** its `MERGED` value is `no` (its work is not yet on `main`)
- **THEN** the output includes the TMUX and MERGED columns

### Requirement: Status separates change spawns from proposal spawns
`jjay status` SHALL group spawns by kind, derived from the name prefix: apply spawns (`app-*`, which track an openspec change) and proposal spawns (`prop-*`, prompt-seeded with no change yet) SHALL be shown in two distinct tables. Change-only columns (MERGED, ARCHIVED, TASKS) apply to the CHANGES table; proposal spawns SHALL NOT be forced into change-shaped columns that are meaningless for them.

#### Scenario: Two tables shown
- **WHEN** both an `app-add-foo` and a `prop-dark-mode` spawn exist and `jjay status` runs
- **THEN** `app-add-foo` appears under a CHANGES table
- **THEN** `prop-dark-mode` appears under a PROPOSAL SPAWNS table

#### Scenario: Only one kind present
- **WHEN** only proposal spawns exist
- **THEN** the PROPOSAL SPAWNS table is shown and the CHANGES table is empty or omitted
