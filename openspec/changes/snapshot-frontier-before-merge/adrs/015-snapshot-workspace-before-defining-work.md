# ADR-015: Snapshot the workspace before defining its work; empty frontier ⇒ unproven

**Status**: Proposed

**Extends**: ADR-013 (verification-gated merge). Composes with ADR-010 (`advanceMainToHead`).

## Context

ADR-013 defined the workspace's work as `WORK = ancestors(<change>@) & main.. & ~empty()` and gates merge success on a smoke test over `WORK`. Two residual escapes were logged:

- **jjay-4oyy (uncommitted).** The work was uncommitted in the spawned workspace's `@`. jj snapshots a workspace's working copy only when a jj command runs *inside that workspace's directory*; `merge` runs from the main session, so the dirty `@` was never snapshotted. `WORK` came out empty → smoke test passed vacuously → merge forgot the workspace and reported "verified" while landing nothing.
- **jjay-ychu (orphan sibling).** Work committed on a sibling `@` never descended from — unreachable from `<change>@`. The pure case (where `@` was moved off the work) also leaves no on-disk trace.

Spike findings (validated against real jj, per ADR-010's lesson):
- `jj -R <wsDir> status` from the main session **snapshots** that workspace's dirty `@` (empty → non-empty; files committed at clean paths). So uncommitted work is recoverable into the frontier cheaply and robustly.
- For the *pure* orphan, after moving `@` off the work the files are also checked out of the tree — so neither the frontier nor a working-dir inspection can see them. Only op-log forensics could, and that is fragile.

## Options Considered

- **Status quo (ADR-013 only).** Empty `WORK` ⇒ "nothing to prove" ⇒ pass + forget. Silently loses uncommitted work; reports false success. Rejected (the bug).
- **`update-stale` / reactive only.** Doesn't snapshot a *clean-but-dirty-@* workspace from the main session; doesn't address the false-success-on-empty.
- **Op-log forensics now (ychu Tier 3).** Could name orphan siblings, but fragile and jj-version-dependent; ADR-013 already scoped it out. Deferred.
- **Snapshot-then-refuse-on-empty (chosen).** Force-snapshot the workspace before computing `WORK` (captures uncommitted work — closes 4oyy); then if `WORK` is still empty, refuse to claim success (keep + warn — saves the orphan case from silent loss, satisfying ychu's acceptance). Cheap, robust, no op-log parsing.

## Decision

`jjay merge` SHALL, before computing the work frontier:

1. **Force-snapshot the target workspace** via `jj -R <wsDir> status` (a snapshotting no-op), so uncommitted `@` work becomes a commit and enters `ancestors(<change>@)`.

Then, in the smoke test:

2. **An empty frontier is "unproven", not "nothing to prove".** If `WORK` is empty after the snapshot, merge SHALL NOT forget the workspace and SHALL NOT report "verified". It keeps the workspace intact, emits a loud warning including the pre-merge recovery handle (`jj op restore <id>`), and exits non-zero.
3. **Dir-content enrichment.** When the workspace directory holds content main lacks (e.g. via `jj -R <wsDir> diff --from main --summary`), the warning SHALL name it. The pure orphan leaves no on-disk trace, so the warning is generic there.

## Consequences

- **Positive**: 4oyy closed — uncommitted work is snapshotted into the frontier, found, merged, and verified.
- **Positive**: ychu's acceptance met — no merge ever reports "verified" on an empty frontier, so the orphan case is a loud, recoverable failure instead of silent loss.
- **Positive**: No fragile op-log parsing; both tiers are `jj -R <wsDir>` invocations.
- **Negative**: A genuinely empty spawn (no work at all) now exits non-zero on merge instead of succeeding vacuously. Acceptable: merging an empty workspace is almost always a mistake; the warning explains it.
- **Negative**: The pure orphan is *detected* (refuse-on-empty) but not *named* — Tier 3 (op-log, jjay-ychu) is still required to identify which sibling commit holds the work.
- **Negative**: `jj -R <wsDir> status` mutates the spawned workspace (creates a snapshot commit) as a side effect of merge — intended, but a write where ADR-013's reads were not.
