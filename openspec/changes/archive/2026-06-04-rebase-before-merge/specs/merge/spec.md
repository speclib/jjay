## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Workspace files preserved after merge
After merge, all files added or modified by the workspace SHALL be present in the merge result. No files SHALL be silently dropped.

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

### Requirement: E2E test scenarios for merge
Integration tests (build tag `integration`) SHALL cover 6 merge scenarios.

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
- **THEN** merge succeeds, new files are present (THE CRITICAL BUG FIX)

#### Scenario 5: Empty workspace
- **WHEN** workspace has no changes
- **THEN** warning printed, merge proceeds

#### Scenario 6: Multiple workspace commits
- **WHEN** workspace has multiple commits
- **THEN** all commits are rebased, merge includes all changes
