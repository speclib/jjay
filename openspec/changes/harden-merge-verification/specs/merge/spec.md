## ADDED Requirements

### Requirement: Merge operates on the workspace's full work frontier
`jjay merge <change>` SHALL define the workspace's work as all non-empty commits reachable from the workspace's heads that are not yet on main — `main..(heads of the workspace's commits) & ~empty()` — not only `<change>@`. The rebase and merge SHALL operate over that frontier, so work committed in `@-` or on a divergent sibling commit is included rather than silently dropped.

#### Scenario: Work in @- while @ is empty
- **WHEN** the workspace's `@` is empty but real work is committed in `@-`
- **THEN** merge includes the `@-` work (the frontier covers `main..@-`)
- **THEN** the merge commit is non-empty and the `@-` work is present on main

#### Scenario: Work on a divergent sibling commit
- **WHEN** the workspace's real work lives on a commit that `<change>@` does not descend from (a sibling)
- **THEN** merge includes that sibling commit's work via the frontier
- **THEN** the sibling work is present on main after merge (not orphaned)

### Requirement: Merge proves the work landed (post-merge smoke test)
After merge, `jjay merge` SHALL verify that the workspace's work actually landed on main before reporting success. (L1) If the workspace had work but main gained no changes, the merge SHALL fail. (L2) Every file added or modified across the work frontier, captured before the merge, SHALL be present on main; any missing file SHALL fail the smoke test. The smoke test SHALL be verbose by default (pre-1.0). Content equivalence (L3) is out of scope for this change.

#### Scenario: Empty merge detected
- **WHEN** the workspace had work but the merge added nothing to main
- **THEN** the smoke test fails (L1) and merge does not report success

#### Scenario: Missing file detected
- **WHEN** a file recorded in the pre-merge work-frontier capture is absent from main after merge
- **THEN** the smoke test fails (L2), naming the expected and missing files

### Requirement: Merge captures a pre-merge recovery handle
Before rewriting any commit, `jjay merge` SHALL capture the current `jj` operation id (from `jj op log`) and SHALL include it in any failure message as a recovery handle (`jj op restore <id>`).

#### Scenario: Recovery handle on failure
- **WHEN** the smoke test fails
- **THEN** the failure message includes the pre-merge operation id and the `jj op restore` command to recover

### Requirement: Workspace lifecycle is gated on merge verification
`jjay merge` SHALL gate the spawned workspace's lifecycle on the smoke test. On success, it SHALL forget the jj workspace so no stale working-copy pointer remains. On failure, it SHALL keep the workspace intact, SHALL NOT advance further, SHALL emit a loud structured warning, and SHALL exit non-zero. Merge SHALL NOT perform any destructive rollback automatically.

#### Scenario: Proven merge forgets the workspace
- **WHEN** the smoke test passes
- **THEN** the jj workspace is forgotten
- **THEN** `jj status` reports no stale working copy for that workspace (there is no live pointer to go stale)

#### Scenario: Unproven merge keeps the workspace non-destructively
- **WHEN** the smoke test fails
- **THEN** the workspace is kept (not forgotten)
- **THEN** the bookmark is not advanced beyond the merge already performed and no rollback is run automatically
- **THEN** merge exits non-zero with a warning naming missing files and the recovery handle

## MODIFIED Requirements

### Requirement: Merge updates main bookmark
After creating the merge commit, the command SHALL move the `main` bookmark to the new merge commit and create a fresh empty change for the user. The bookmark SHALL only be considered final once the post-merge smoke test passes; on smoke-test failure the command SHALL stop without forgetting the workspace, leaving the state inspectable for recovery.

#### Scenario: Bookmark moved
- **WHEN** `jjay merge feat-payments` completes successfully and the smoke test passes
- **THEN** `main` bookmark points to the merge commit
- **THEN** the user's working copy is a fresh empty change on top of main
- **THEN** the merged workspace is forgotten

### Requirement: E2E test scenarios for merge
Integration tests (build tag `integration`) SHALL cover the merge scenarios below, including the three reproduced jjay-q6ko failure topologies (empty `@` with work in `@-`, divergent-sibling work, and post-merge workspace staleness/forget), plus an unproven-keeps-workspace scenario and an L1 empty-land detector.

#### Scenario 1: Clean merge — no main changes
- **WHEN** main has no new commits since workspace was spawned
- **THEN** merge succeeds, all workspace files present

#### Scenario 2: Main moved forward, no overlap
- **WHEN** main and workspace modify different files
- **THEN** merge succeeds, files from both sides present

#### Scenario 3: Main moved forward, same file modified
- **WHEN** main and workspace both modify the same file
- **THEN** rebase surfaces conflict, merge aborts with clear message

#### Scenario 4: Workspace adds new files
- **WHEN** workspace adds files that main doesn't have
- **THEN** merge succeeds, new files are present

#### Scenario 5: Empty workspace
- **WHEN** workspace has no changes
- **THEN** warning printed, merge proceeds

#### Scenario 6: Multiple workspace commits
- **WHEN** workspace has multiple commits
- **THEN** all commits are rebased, merge includes all changes

#### Scenario 7: Main adds new files ahead of bookmark (TestMerge_MainAddsNewFiles)
- **WHEN** main commits a new file/directory ahead of the `main` bookmark after the workspace base, and the workspace commits unrelated work
- **THEN** after merge both the main-side addition and the workspace work exist on `main`

#### Scenario 8: Empty @ with work in @- (TestMerge_EmptyAtWorkInParent)
- **WHEN** the workspace's `@` is empty and its real work is in `@-`
- **THEN** merge lands the `@-` work on main (non-empty merge), proving the frontier reaches `@-`

#### Scenario 9: Work on a divergent sibling (TestMerge_WorkOnSibling)
- **WHEN** the workspace's work is on a sibling commit that `<change>@` does not descend from
- **THEN** merge lands the sibling work on main (it is not orphaned)

#### Scenario 10: Workspace not stale after a proven merge (TestMerge_NotStaleAfterMerge)
- **WHEN** a normal merge succeeds and its smoke test passes
- **THEN** the workspace is forgotten so no stale working-copy pointer remains

#### Scenario 11: Unproven merge keeps the workspace (TestMerge_UnprovenKeepsWorkspace)
- **WHEN** the smoke test detects missing/empty-landed work
- **THEN** the workspace is kept, merge exits non-zero, and the warning includes the recovery handle
