### Requirement: Merge creates merge commit
The `jjay merge <change>` command SHALL first rebase the workspace branch onto current main via `jj rebase -b <change>@ -d main`, then create a merge commit. If the rebase surfaces conflicts, the command SHALL report them and exit without merging. The merge commit message SHALL be `merge <change> into main`.

#### Scenario: Clean rebase and merge
- **WHEN** `jjay merge feat-payments` is executed and the workspace has no conflicts with main
- **THEN** `jj rebase -b feat-payments@ -d main` is run first
- **THEN** the merge commit is created with main and the rebased workspace as parents
- **THEN** all files from both main and workspace are present

#### Scenario: Rebase with conflicts
- **WHEN** `jjay merge feat-payments` is executed and rebase produces conflicts
- **THEN** jjay reports the conflicts
- **THEN** jjay exits with non-zero exit code
- **THEN** no merge commit is created
- **THEN** the user can resolve conflicts manually and retry

### Requirement: Merge updates main bookmark
After creating the merge commit, the command SHALL move the `main` bookmark to the new merge commit and create a fresh empty change for the user. The bookmark SHALL only be considered final once the post-merge smoke test passes; on smoke-test failure the command SHALL stop without forgetting the workspace, leaving the state inspectable for recovery.

#### Scenario: Bookmark moved
- **WHEN** `jjay merge feat-payments` completes successfully and the smoke test passes
- **THEN** `main` bookmark points to the merge commit
- **THEN** the user's working copy is a fresh empty change on top of main
- **THEN** the merged workspace is forgotten

### Requirement: Merge requires workspace to exist
The command SHALL verify the jj workspace exists before proceeding.

#### Scenario: Workspace exists
- **WHEN** `jjay merge feat-payments` is executed and workspace `feat-payments` exists
- **THEN** merge proceeds

#### Scenario: Workspace does not exist
- **WHEN** `jjay merge feat-payments` is executed and workspace `feat-payments` does not exist
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the workspace does not exist

### Requirement: Merge warns on empty workspace
The command SHALL warn if the workspace's working copy is empty (no changes). It SHALL still proceed after the warning.

#### Scenario: Empty workspace
- **WHEN** `jjay merge feat-payments` is executed and the workspace's `@` is empty
- **THEN** a warning is printed indicating the workspace has no changes
- **THEN** merge proceeds anyway

### Requirement: Merge requires change name argument
The command SHALL require exactly one argument: the change name.

#### Scenario: No argument
- **WHEN** `jjay merge` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code

### Requirement: Merge does not push
The command SHALL NOT push to any remote. Pushing is a separate user action.

#### Scenario: No push
- **WHEN** `jjay merge feat-payments` completes successfully
- **THEN** no `jj git push` is executed
- **THEN** the user can push manually when ready

### Requirement: Workspace files preserved after merge
After merge, all files added or modified by the workspace SHALL be present in the merge result, AND all committed work on the main line SHALL be preserved — including commits in the main working copy (`@`) that are **ahead of the `main` bookmark** at merge time. Because jj auto-snapshots the working copy, uncommitted main-side edits are captured into `@` and are therefore also preserved. No files SHALL be silently dropped from either side.

#### Scenario: Workspace adds new files, main adds different files
- **WHEN** main adds `bar.txt` and workspace adds `baz.txt`
- **THEN** after merge, both `bar.txt` and `baz.txt` exist

#### Scenario: Workspace modifies file also modified on main
- **WHEN** main modifies `foo.txt` and workspace also modifies `foo.txt`
- **THEN** after rebase, jj surfaces the conflict explicitly
- **THEN** merge does not proceed until conflict is resolved

#### Scenario: Workspace checks off tasks in tasks.md
- **WHEN** main has `tasks.md` with unchecked boxes and workspace checks them off
- **THEN** after merge, the checked-off version is present (rebase applies workspace's changes on top of main)

#### Scenario: Main work ahead of bookmark survives merge
- **WHEN** a workspace `feat` is spawned, then a new directory `openspec/changes/new-thing/` is committed in the main working copy ahead of the `main` bookmark, and `jjay merge feat` runs
- **THEN** after the merge `openspec/changes/new-thing/` exists on `main`
- **THEN** the workspace's own work also exists on `main`

#### Scenario: Uncommitted main edits are snapshotted and preserved
- **WHEN** the main working copy has uncommitted edits at merge time
- **THEN** jj snapshots them into `@`, so they are ahead-of-bookmark work and are folded into the merge like any committed main work
- **THEN** the edits are present on `main` after merge (jj has no unreachable "dirty" state that would require aborting)

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

#### Scenario 9: Divergent sibling work is detected by the smoke test (TestMerge_SiblingDetected)
- **WHEN** the workspace's work is on a sibling commit that `<change>@` does not descend from
- **THEN** merge does not auto-include it (unreachable from `<change>@`)
- **THEN** the smoke test fails loudly, the workspace is kept, and a recovery handle is given (not silent success)

#### Scenario 10: Workspace not stale after a proven merge (TestMerge_NotStaleAfterMerge)
- **WHEN** a normal merge succeeds and its smoke test passes
- **THEN** the workspace is forgotten so no stale working-copy pointer remains

#### Scenario 11: Unproven merge keeps the workspace (TestMerge_UnprovenKeepsWorkspace)
- **WHEN** the smoke test detects missing/empty-landed work
- **THEN** the workspace is kept, merge exits non-zero, and the warning includes the recovery handle

### Requirement: Merge operates on the workspace's ancestor work frontier
`jjay merge <change>` SHALL define the workspace's work as all non-empty commits in `ancestors(<change>@) & main.. & ~empty()` — i.e. everything on `<change>@`'s line not yet on main, including `@-` — not only `<change>@` itself. The rebase and merge SHALL operate over that frontier, so work committed in `@-` (or any ancestor of `@`) is included rather than silently dropped.

NOTE — limitation (spike-confirmed): jj associates exactly one working-copy commit (`@`) with a workspace, so a commit the workspace's `@` never descended from (a true divergent *sibling*, e.g. the `zroyto` orphan in q6ko instance 3) is **not reachable from `<change>@`** and cannot be auto-included. Such work is handled by the post-merge smoke test below, which fails loudly and keeps the workspace for recovery rather than silently dropping it. Auto-including orphan siblings would require op-log scanning and is a documented non-goal.

#### Scenario: Work in @- while @ is empty
- **WHEN** the workspace's `@` is empty but real work is committed in `@-`
- **THEN** merge includes the `@-` work (the frontier covers `main..@-`)
- **THEN** the merge commit is non-empty and the `@-` work is present on main

#### Scenario: Divergent sibling work is detected, not silently dropped
- **WHEN** the workspace's real work lives on a sibling commit that `<change>@` does not descend from
- **THEN** the work is not auto-included (it is unreachable from `<change>@`)
- **THEN** the post-merge smoke test fails loudly, the workspace is kept, and the failure names the recovery handle — instead of reporting false success

### Requirement: Merge proves the work landed (post-merge smoke test)
After merge, `jjay merge` SHALL verify that the workspace's work actually landed on main before reporting success. (L1) If the workspace had work but main gained no changes, the merge SHALL fail. (L2) Every file **added or modified** across the work frontier, captured before the merge, SHALL be present on main; any missing file SHALL fail the smoke test. The capture SHALL exclude files the workspace **deleted** (a net delete is legitimately absent from main), and L2 presence SHALL be satisfied by either the exact path OR the file's basename appearing on main, so that a file legitimately **moved/renamed** by the merge — most importantly an `/opsx:archive` that moves `openspec/changes/X/` → `openspec/changes/archive/<date>-X/` — is not falsely reported missing. The smoke test SHALL be verbose by default (pre-1.0). Content equivalence (L3) is out of scope for this change.

#### Scenario: Empty merge detected
- **WHEN** the workspace had work but the merge added nothing to main
- **THEN** the smoke test fails (L1) and merge does not report success

#### Scenario: Missing file detected
- **WHEN** a file recorded in the pre-merge work-frontier capture is absent from main after merge (by neither its path nor its basename)
- **THEN** the smoke test fails (L2), naming the expected and missing files

#### Scenario: Archive-move is not a false positive
- **WHEN** the merge archives the change, moving a captured file from `openspec/changes/X/` to `openspec/changes/archive/<date>-X/` (same basename, different directory)
- **THEN** L2 treats the file as present (matched by basename) and the smoke test passes

#### Scenario: Deleted file is not expected on main
- **WHEN** the workspace's net effect deletes a file
- **THEN** that file is excluded from the captured set and its absence from main does not fail the smoke test

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
