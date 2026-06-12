## Why

`jjay merge` (ADR-013) defines the workspace's work as `ancestors(<change>@) & main.. & ~empty()` and gates success on a smoke test. Two ways the *real* work escapes that frontier remain:

- **[jjay-4oyy](../../.beans/jjay-4oyy--merge-smoke-test-misses-uncommitted-workspace-work.md) (uncommitted, HIGH — happened).** The work sat **uncommitted** in the spawned workspace's `@`. jj auto-snapshots only when a command runs *inside* that workspace dir; `merge` runs from the main session, so the dirty `@` was never snapshotted. The frontier was empty → smoke test "nothing to prove" → merge forgot the workspace and reported **"verified"** while landing nothing. (`prop-no-respawn-claude-apply`: a full proposal lost from jj's view; recovered only because `jj workspace forget` leaves the dir on disk.)
- **[jjay-ychu](../../.beans/jjay-ychu--detect-orphaned-workspace-commits-via-jj-op-log-tr.md) (orphan sibling).** Work committed on a sibling `@` never descended from. Unreachable from `<change>@`; the pure case (where `@` was moved off the work) leaves **no on-disk trace** either.

Both are escape routes from the same trap: **merge only looks where `<change>@`'s ancestry points.** This change closes 4oyy and satisfies ychu's *acceptance* ("detect & warn instead of a silent empty merge") without the fragile op-log forensics — which stays deferred in ychu.

**Spike-verified:** `jj -R <wsDir> status` from the main session **does** snapshot a sibling workspace's dirty `@` (empty → non-empty, files committed at clean paths). So the fix is robust, not fragile.

## What Changes

Two tiers, both cheap and robust, added to `jjay merge`'s pre-flight (composing with ADR-010 `advanceMainToHead` and ADR-013's smoke test):

- **Tier 1 — force-snapshot the workspace before computing the frontier.** Run `jj -R <wsDir> status` (a harmless, snapshotting no-op) for the target workspace *before* `WORK` is computed, so uncommitted `@` work becomes a real commit and enters the frontier. Closes 4oyy: the work is then found, merged, and verified normally.
- **Tier 2 — refuse to claim success on an empty frontier (the guarantee).** After the snapshot, if `WORK` is **still empty**, merge SHALL NOT forget the workspace and SHALL NOT report "verified": it keeps the workspace, emits a loud warning with the recovery handle (`jj op restore <preMergeOp>`), and exits non-zero. This is what saves the pure-orphan case (ychu): **no silent empty merge.**
  - **Dir-check enrichment.** When the workspace directory *does* still hold content main lacks (e.g. an `openspec/changes/<x>/` the snapshot didn't surface into the frontier), the warning SHALL name it ("workspace dir has X, not on main") so the human/AI knows what to recover. (The *pure* orphan leaves no on-disk trace, so the warning is generic there — Tier 3 / op-log would be needed to name it.)

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `merge`: before computing the work frontier, merge SHALL force-snapshot the target workspace (`jj -R <wsDir> status`) so uncommitted `@` work is captured; and SHALL treat a post-snapshot empty frontier as **unproven** — keeping the workspace, warning loudly (with a dir-content hint when available and the recovery handle), and exiting non-zero rather than forgetting it and reporting success.

## Impact

- **Code**: `internal/merge/merge.go` — a pre-frontier snapshot step (`jj -R <wsDir> status`); the empty-frontier branch of the smoke test becomes "unproven, keep + warn" instead of "nothing to prove, pass"; an optional `jj -R <wsDir> diff --from main --summary` for the dir-content hint. Reuses the existing `preMergeOp` recovery handle. New integration tests for the uncommitted-`@` and empty-frontier cases.
- **Severity**: HIGH — closes a "reports verified while losing work" footgun on the no-AI core (4oyy).
- **Relation**: extends ADR-013 (`harden-merge-verification`); composes with ADR-010. Closes **jjay-4oyy**; satisfies **jjay-ychu**'s detection acceptance (ychu remains open for the deferred Tier-3 op-log *identification* of orphan siblings).
- **ADRs**: ADR-015 (snapshot the workspace before defining its work; empty frontier ⇒ unproven, never silent success).
- **Beans**: jjay-4oyy → in-progress, linked here; jjay-ychu → note that detection is satisfied here, Tier-3 remains its scope.

## Deferred / Out of Scope

- **Tier 3 — op-log forensics (jjay-ychu).** Actually *identifying and naming* an orphaned sibling commit (`@` never visited, no on-disk trace) needs scanning `jj op log` for commits authored in the workspace — fragile, jj-version-dependent. Deferred; only needed if Tier 2's refuse-and-warn proves insufficient in practice.
- **Auto-recovery / auto-rollback.** Merge stays non-destructive: it keeps the workspace and hands the human a recovery handle. No auto-include of hidden work.
