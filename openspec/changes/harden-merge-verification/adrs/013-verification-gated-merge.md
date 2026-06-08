# ADR-013: Verification-gated merge — prove the work landed, gate workspace lifecycle on proof

**Status**: Proposed

## Context

`jjay merge` (`internal/merge/merge.go`) operates on a single revset, `<change>@`, for the rebase, the merge commit, and all checks. This assumes `<change>@` is the workspace's work. Bean jjay-q6ko reproduced three failures of that assumption:

1. **Staleness** — merge rewrites `<change>@` from the main session; the spawned workspace's working-copy pointer is never updated, so `jj status` inside it errors `The working copy is stale`. On essentially every merge.
2. **Empty `@`** — work was in `@-` while `@` was empty; `jj new main <change>@` produced an empty merge commit, reported as success.
3. **Orphaned sibling** — the real proposal lived on a divergent commit `<change>@` did not descend from; merge never saw it and silently landed nothing.

Today merge proves only a negative ("no conflicts"). Conflict-free ≠ merged (jjay-rse4). ADR-010 (`advanceMainToHead`) and `rebase-before-merge` already protect main-side work; nothing protects workspace-side work or surfaces the staleness merge leaves behind.

## Options Considered

- **Status quo (`<change>@`-only, no verify).** Silently merges empty/orphaned work; leaves workspace stale. Rejected (the bug).
- **Forget the workspace unconditionally after merge.** Kills staleness, but if the merge landed nothing (instances 2/3) it destroys the only copy of the work. Rejected — destructive on the exact failure it can't detect.
- **`update-stale` the workspace after merge.** Fixes the stale pointer reactively but still reports success on empty/orphaned merges. Treats the symptom, not the cause.
- **Verification-gated lifecycle (chosen).** Define the work robustly, prove it landed, then forget on success / keep+warn on failure. Closes all three instances with one mechanism and matches the stated stance: keep the workspace until the merge is proven successful.

## Decision

`jjay merge` SHALL be a verification-gated pipeline:

1. **Snapshot** the `jj op log` head before any rewrite (recovery handle).
2. **Define the work** as `WORK = main..(heads of the workspace's commits) & ~empty()` — including `@-` and divergent siblings, not just `<change>@`. Capture the added/modified file set across `WORK`.
3. **Fold + merge** via the existing `advanceMainToHead` (ADR-010) then rebase/merge over `WORK`.
4. **Smoke test** (rse4 L1+L2): if the workspace had work, main MUST have gained changes (L1); every captured file MUST be present on main (L2).
5. **Gate the lifecycle on the proof:**
   - **Proven** → forget the jj workspace. No live pointer ⇒ staleness is structurally impossible.
   - **Unproven** → keep the workspace intact, do not advance further, emit a loud structured warning (expected vs missing files + `jj op restore <preMergeOp>`), exit non-zero. No auto-rollback.

L3 (content equivalence) is deferred (rebase-combination caveat). Auto-rollback / self-healing are deferred.

## Consequences

- **Positive**: Empty-`@` and orphaned-sibling merges (instances 2/3) can no longer report false success — the smoke test fails loudly with a recovery handle.
- **Positive**: Staleness (instance 1) is closed structurally — a proven merge forgets the workspace, so there is no pointer to go stale; an unproven merge deliberately keeps it.
- **Positive**: Non-destructive on failure — the work is never thrown away on the case the tool can't yet verify.
- **Negative**: Merge gains real logic (frontier revset, file capture, smoke test) — more than "rebase onto bookmark".
- **Negative**: Behavior change — a clean merge now forgets the jj workspace (directory cleanup stays in `cleanup`); must be documented.
- **Negative**: The `heads()` frontier revset must be validated end-to-end against real topologies (ADR-010's lesson), not by unit reasoning alone.
