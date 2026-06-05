## ADDED Requirements

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
