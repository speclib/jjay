## MODIFIED Requirements

### Requirement: Merge operates on the workspace's ancestor work frontier
`jjay merge <change>` SHALL define the workspace's work as all non-empty commits in `ancestors(<change>@) & main.. & ~empty()`. **Before computing this frontier, merge SHALL force-snapshot the target workspace** by running a snapshotting jj command scoped to the workspace directory (`jj -R <wsDir> status`), so work left **uncommitted** in the workspace's `@` is captured into a commit and enters the frontier — `merge` runs from the main session, where jj would otherwise never snapshot the spawned workspace's dirty working copy. A divergent *sibling* commit the workspace `@` never descended from is still unreachable from `<change>@` and is handled by the smoke test (the empty-frontier guard below), not by this frontier.

#### Scenario: Uncommitted @ work is snapshotted into the frontier
- **WHEN** a spawned workspace has real work left uncommitted in its `@` (never snapshotted, because merge runs from the main session) and `jjay merge` runs
- **THEN** merge first force-snapshots the workspace (`jj -R <wsDir> status`)
- **THEN** the now-committed work appears in the frontier and is merged onto main (not silently dropped)

#### Scenario: Work in @- while @ is empty
- **WHEN** the workspace's `@` is empty but real work is committed in `@-`
- **THEN** merge includes the `@-` work (the frontier covers `main..@-`)
- **THEN** the merge commit is non-empty and the `@-` work is present on main

### Requirement: Merge proves the work landed (post-merge smoke test)
After merge, `jjay merge` SHALL verify that the workspace's work actually landed on main before reporting success. (L1) If the workspace had work but main gained no changes, the merge SHALL fail. (L2) Every file added or modified across the work frontier, captured before the merge, SHALL be present on main (matched by path or basename; net-deletes excluded); any missing file SHALL fail the smoke test.

**An empty work frontier (after the pre-merge force-snapshot) SHALL be treated as UNPROVEN, not as "nothing to prove".** When the frontier is empty, merge SHALL NOT forget the workspace and SHALL NOT report "verified": it SHALL keep the workspace intact, emit a loud warning including the recovery handle (`jj op restore <preMergeOp>`), and exit non-zero. When the workspace directory holds content main lacks (e.g. detected via `jj -R <wsDir> diff --from main --summary`), the warning SHALL name that content; for a pure orphan that leaves no on-disk trace, the warning MAY be generic. The smoke test SHALL be verbose by default (pre-1.0). Content equivalence (L3) is out of scope.

#### Scenario: Empty frontier after snapshot is unproven, not success
- **WHEN** the work frontier is empty even after the pre-merge force-snapshot (e.g. a genuinely empty spawn, or work orphaned on an unreachable sibling)
- **THEN** merge does NOT report "verified" and does NOT forget the workspace
- **THEN** merge keeps the workspace, exits non-zero, and the warning includes the recovery handle

#### Scenario: Orphaned sibling work is not silently merged (jjay-ychu acceptance)
- **WHEN** the workspace's real work is on a sibling commit `<change>@` never descended from (empty frontier even after snapshot)
- **THEN** merge refuses to claim success and keeps the workspace for recovery — instead of reporting a verified, empty merge

#### Scenario: Dir-content hint when work is still on disk
- **WHEN** the frontier is empty but the workspace directory still contains files main lacks (e.g. an `openspec/changes/<x>/` directory)
- **THEN** the warning names that content so the recovery target is clear

#### Scenario: Missing file detected
- **WHEN** a file recorded in the pre-merge work-frontier capture is absent from main after merge (by neither path nor basename)
- **THEN** the smoke test fails (L2), naming the expected and missing files
