# ADR-007: Rebase workspace onto main before merging

**Status**: Proposed

## Context

`jjay merge` creates a merge commit with `jj new main <workspace>@`. When both sides have modified the same files, or when the workspace adds files that main doesn't know about, jj's 3-way merge silently resolves by picking one side. This has caused repeated data loss: task checkboxes reverted, blog posts dropped, a bug-tracking bean deleted.

The core issue: without rebase, the merge is a true 3-way merge where jj must reconcile divergent branches. With files added only on one side, jj may interpret the absence on the other side as "file should not exist."

## Options Considered

- **Rebase then merge** — rebase workspace onto main first, making the merge trivial. Conflicts surface explicitly during rebase. Most reliable.
- **Fast-forward only** — rebase then move bookmark without merge commit. Simpler but loses the "merged from workspace X" marker in history.
- **Manual conflict resolution after merge** — keep current approach, add post-merge checks. Doesn't fix the root cause — files are already dropped.
- **Copy files instead of jj merge** — nuclear option, bypasses jj entirely. Loses all VCS benefits.

## Decision

Rebase workspace branch onto current main (`jj rebase -b <ws>@ -d main`) before creating the merge commit. After rebase, the workspace includes all of main's changes, making the merge trivial. If rebase surfaces conflicts, abort and let the user resolve.

Keep the merge commit (not fast-forward) for clear history markers.

## Consequences

- **Positive**: Eliminates silent file drops — the most critical workflow bug
- **Positive**: Real conflicts surface explicitly as jj conflict markers during rebase
- **Positive**: Merge becomes trivial after rebase — no 3-way resolution needed
- **Negative**: Rebase rewrites workspace commit IDs (normal jj behavior, no practical impact)
- **Negative**: If workspace has long-lived branches with many commits, rebase takes longer (negligible for typical agent sessions)
