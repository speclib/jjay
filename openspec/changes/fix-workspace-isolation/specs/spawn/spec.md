## MODIFIED Requirements

### Requirement: Spawn creates jj workspace
The `jjay spawn <change>` command SHALL first run `jj new` in the main workspace to snapshot uncommitted work, then create a jj workspace named `<change>` at `../<project-name>-workspaces/<change>` via `jj workspace add --name <change> --revision @- <path>`. The `@-` revision contains the snapshotted files. The parent directory SHALL be created if it does not exist.

#### Scenario: Workspace created with isolation
- **WHEN** `jjay spawn feat-payments` is executed in project `jjay`
- **THEN** `jj new` is run first in the main workspace
- **THEN** `jj workspace add --name feat-payments --revision @- ../jjay-workspaces/feat-payments` is run
- **THEN** the child workspace contains all files from the snapshot
- **THEN** the main workspace's `@` is a fresh empty change (nothing to lose)

#### Scenario: Main workspace safe during concurrent work
- **WHEN** a spawned agent creates jj operations in the child workspace
- **THEN** the main workspace may become stale
- **THEN** running `jj workspace update-stale` in the main workspace does NOT lose work (because `@` is empty, all prior work is in `@-`)

#### Scenario: Workspace already exists
- **WHEN** `jjay spawn feat-payments` is executed and jj workspace `feat-payments` already exists
- **THEN** jjay exits with a non-zero exit code
- **THEN** `jj new` has NOT been run (preconditions checked before any mutations)
- **THEN** an error message is printed indicating the workspace already exists
