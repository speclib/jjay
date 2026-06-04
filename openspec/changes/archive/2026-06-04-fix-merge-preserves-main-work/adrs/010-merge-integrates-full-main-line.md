# ADR-010: Merge integrates the full main line, not only the `main` bookmark

**Status**: Accepted

## Context

`jjay merge` (`internal/merge/merge.go`) rebases the workspace onto the `main` bookmark, creates a merge commit `jj new main <change>@`, then advances the bookmark to it. This assumes the `main` bookmark is the tip of all main-side work. It often isn't: in normal jjay use, the orchestrator does work in the main **working copy** (`@`) — creating proposals, editing beans — that lands in commits *ahead of* the `main` bookmark and is not folded back into it until something explicitly moves the bookmark. When merge then bases the merge commit on the lagging bookmark, that ahead-of-bookmark work is in neither merge parent and is silently dropped when the bookmark advances to the merge commit.

This is distinct from the `rebase-before-merge` fix (ADR for jjay-30gc), which addressed files modified on *both* sides via 3-way silent picking. Here nothing conflicts — the main-side work simply isn't part of what merge considers "main".

## Options Considered

- **Status quo (bookmark-only)** — base merge on the `main` bookmark. Silently loses ahead-of-bookmark main work. Rejected (the bug).
- **Advance the bookmark to `@` before merging** — fold the main working copy into `main` first, then merge onto it. Simple and faithful to "main is everything committed on the main line." Risk: requires `@`'s main-side work to be committed (not a dirty working copy).
- **Rebase the workspace onto the head of main's working-copy line** (rather than the bookmark) and merge there. Equivalent effect; chooses the true tip as the destination.
- **Detect-and-abort only** — refuse to merge when `@` is ahead of the bookmark, telling the user to advance main first. Safe but pushes manual work onto the user every merge.

## Decision

`jjay merge` SHALL treat the **tip of the main working-copy line**, not the `main` bookmark alone, as the integration target:

- Before merging, **detect whether `@` (main working copy) is ahead of the `main` bookmark**. If so, the ahead-of-bookmark commits are part of "main" and MUST be included in the merge.
- **Include them** by advancing the `main` bookmark to the true main tip (`latest(main..@ & ~empty())`) before the rebase/merge, so the merge commit's tree is `full-main-line ∪ change@`.
- **No dirty-state abort.** Implementation found jj auto-snapshots the working copy, so uncommitted main edits are captured into `@` and folded in like any committed work — there is no unreachable dirty state requiring an abort (this differs from a git-based design). The existing rebase-conflict abort is preserved for genuine content conflicts.

## Consequences

- **Positive**: Main-side work created after a spawn (new proposals, bean edits) survives merge — the data-loss class in jjay-ug7y is closed.
- **Positive**: Failure is loud and non-destructive; the bookmark only moves when the full main line is integrated.
- **Negative**: Merge must reason about bookmark-vs-working-copy divergence — more logic than "rebase onto bookmark".
- **Negative**: If main's working copy is dirty, merge may now refuse where it previously (wrongly) "succeeded" by losing data — a behavior change, but the correct one.
- **Negative**: Interacts with the snapshot `jj new` that spawn performs; the fix must be validated against the real spawn→work-on-main→merge lifecycle, not just unit-level rebase tests.
