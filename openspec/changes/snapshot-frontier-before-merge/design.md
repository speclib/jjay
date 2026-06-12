## Context

`internal/merge/merge.go` (ADR-013) computes `WORK = ancestors(<change>@) & main.. & ~empty()`, captures the changed-file set, merges, then runs an L1/L2 smoke test. On a clean pass it forgets the workspace; on failure it keeps it + warns. Two residual escapes (jjay-4oyy, jjay-ychu) share one root: **merge only sees `ancestors(<change>@)`'s committed work.** 4oyy's work was uncommitted in `@` (never snapshotted, because merge runs from the main session); ychu's was on an unreachable sibling.

Spikes (validated against real jj):
- `jj -R <wsDir> status` from the main session snapshots that workspace's dirty `@` → uncommitted work becomes committed and reachable. Robust, no parsing.
- A *pure* orphan (`@` moved off the work) checks the work out of the tree too — so neither the frontier nor a dir inspection sees it; only op-log (deferred) could name it.

## Goals / Non-Goals

**Goals:**
- Snapshot the workspace before defining `WORK` so uncommitted work is captured (close 4oyy).
- Never report "verified" on an empty frontier — keep + warn instead (satisfy ychu's no-silent-empty-merge acceptance).
- Name on-disk content in the warning when present (enrichment).

**Non-Goals:**
- Op-log forensics to *identify* an orphan sibling (jjay-ychu Tier 3) — deferred.
- Auto-recovery / auto-rollback — merge stays non-destructive, hands back a recovery handle.
- Changing the L1/L2 logic for the non-empty case (ADR-013 stands).

## Decisions

- **Pre-frontier snapshot.** In `Merge`, before `workFrontierFiles`, resolve the workspace dir and run `jj -R <wsDir> status` (ignore its stdout; it snapshots as a side effect). Order: after `advanceMainToHead` (ADR-010), before frontier capture. The workspace dir is the same `WorkspaceDir(name, root)` merge already knows.
- **Empty frontier ⇒ unproven.** Today `smokeTest` returns nil when `expectedFiles` is empty ("nothing to prove"). Change: an empty frontier (post-snapshot) is a failure — return a loud error (kept workspace, recovery handle), so the caller does NOT reach the forget / `jj new` / "verified" path. The existing `checkWorkspaceEmpty` warning stays as an early heads-up.
- **Dir-content hint.** Before erroring on empty frontier, run `jj -R <wsDir> diff --from main --summary`; if non-empty, include those paths in the warning. Empty (pure orphan) ⇒ generic message. Best-effort: a failure to compute the hint doesn't change the refuse-on-empty outcome.
- **Reuse `preMergeOp`.** The recovery handle already captured in ADR-013 is threaded into the new warning.

## Risks / Trade-offs

- **`jj -R <wsDir> status` mutates the spawned workspace** (creates a snapshot commit). Intended, but it is a write in merge's pre-flight where ADR-013 only read. Validate it doesn't disturb the main session's `@`.
- **Genuinely-empty spawn now fails merge.** A spawn with no work at all exits non-zero instead of vacuously succeeding. Acceptable (merging empty is a mistake); the warning explains it. Note: this is a behavior change for any caller/test that merged an empty workspace expecting success — audit `TestMerge_EmptyWorkspace`.
- **Pure orphan still only detected, not named.** Honest limitation; Tier 3 (op-log) remains in ychu.
- **Snapshot from main session edge cases** — verify `jj -R <wsDir>` works when the workspace is stale (q6ko territory) or when no tmux/agent is attached. Spike covered the clean dirty-@ case; the implementer must validate stale/edge cases end-to-end (ADR-010 lesson).
