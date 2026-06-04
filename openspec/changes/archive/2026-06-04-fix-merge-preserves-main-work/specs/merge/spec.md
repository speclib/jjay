## MODIFIED Requirements

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
Integration tests (build tag `integration`) SHALL cover the merge scenarios below, including a scenario for main-side work created after the workspace base (the mirror of the workspace-adds-new-files scenario).

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
