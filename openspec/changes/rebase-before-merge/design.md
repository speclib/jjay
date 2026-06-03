## Context

`jjay merge` creates a merge commit via `jj new main <workspace>@`. When both sides modify the same files or when the workspace adds new files, jj's 3-way merge silently picks one side. This has caused repeated data loss: task checkboxes reverted, blog posts dropped, beans deleted.

## Goals / Non-Goals

**Goals:**
- Eliminate silent file drops by rebasing before merging
- Surface real conflicts explicitly so the user can resolve them
- 6 e2e test scenarios proving merge is rock solid

**Non-Goals:**
- Automatic conflict resolution
- Changing the merge commit format (still a two-parent merge)
- Fast-forward merges (keep merge commits for clear history)

## Decisions

### Rebase workspace onto main before merge

New merge sequence:

```
Before:
  main: A ── B ── C
          \
  ws:      └── D ── E

Step 1: jj rebase -b <ws>@ -d main
  main: A ── B ── C
                    \
  ws (rebased):      └── D' ── E'

Step 2: jj new main <ws>@ -m "merge <ws> into main"
  main: A ── B ── C ──────── M  (trivial merge, E' already has C)
                    \       /
  ws:                D' ── E'

Step 3: jj bookmark set main -r @
Step 4: jj new
```

After rebase, the workspace includes all of main's changes. The merge commit is trivial — no 3-way conflict resolution needed.

_Alternative: fast-forward main to workspace (skip merge commit) — rejected because merge commits provide clear "this came from workspace X" markers in history._

### Check for conflicts after rebase

After `jj rebase`, check if the rebased commits have conflicts:

```bash
jj log -r "<ws>@" --no-graph -T 'if(conflict, "conflict", "clean")'
```

If conflicts exist, report them and exit. The user resolves manually in the workspace, then retries `jjay merge`.

### Keep merge commit (not fast-forward)

Even though the rebase makes the merge trivial, we still create a merge commit. This gives a clear marker in history: "this workspace's work was merged into main at this point."

### E2E test structure

Each scenario:
1. Create temp jj repo
2. Set up main and workspace with specific file patterns
3. Run `jjay merge`
4. Assert file presence/absence in the result

Tests use `//go:build integration` tag. Each scenario is a separate `t.Run()` subtest.

```go
func TestMerge_CleanNoMainChanges(t *testing.T) { ... }
func TestMerge_MainMovedNoOverlap(t *testing.T) { ... }
func TestMerge_SameFileModified(t *testing.T) { ... }
func TestMerge_WorkspaceAddsNewFiles(t *testing.T) { ... }
func TestMerge_EmptyWorkspace(t *testing.T) { ... }
func TestMerge_MultipleWorkspaceCommits(t *testing.T) { ... }
```

## Risks / Trade-offs

- [Rebase rewrites workspace history] → Acceptable. jj handles this naturally. The workspace's commits get new IDs but the content is preserved.
- [Rebase may fail on complex histories] → Unlikely for single-workspace branches. If it happens, the error is clear.
- [Extra step adds latency] → Negligible. Rebase is fast for small branches.
