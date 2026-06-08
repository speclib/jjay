## Why

`jjay merge` "almost always gives a stale error", and worse, sometimes reports success while landing **nothing**. Bean [jjay-q6ko](../../.beans/jjay-q6ko--merging-now-almost-always-give-a-stale-error.md) logs three reproduced instances of the same root assumption breaking: **`merge` treats `<change>@` as "the workspace's work".** It isn't always.

- **Instance 1 — staleness.** Merge rebases `<change>@` and builds the merge commit from the **main** session, rewriting `<change>@` out from under the spawned workspace. Nothing updates that workspace's working-copy pointer, so `jj status` inside it errors `The working copy is stale`. Happens on essentially every merge, because every merge rewrites `<change>@`.
- **Instance 2 — empty `@`.** The real work was in `@-`; `@` was empty. `jj new main <change>@` built an **empty merge commit** ("0 files changed"). Merge reported success; nothing visibly changed.
- **Instance 3 — orphaned sibling.** The actual proposal lived on a **divergent** commit (`zroyto`) that `<change>@` did not descend from. `merge` only ever follows `<change>@`'s line, so it never saw it. Merge reported success; the real proposal was left orphaned off-main entirely.

This is the sibling of [jjay-ug7y](../../.beans/jjay-ug7y--critical-merge-deletes-active-sibling-change-dirs.md) (fixed by ADR-010, `advanceMainToHead`) and [jjay-30gc](../../.beans/jjay-30gc--critical-openspec-artifacts-diverge-on-merge-files.md) (fixed by `rebase-before-merge`). Those closed *main-side* loss. This closes *workspace-side* loss and the staleness it leaves behind. It **absorbs** the fully-specced smoke-test bean [jjay-rse4](../../.beans/jjay-rse4--post-merge-smoke-test-prove-the-merge-actually-lan.md) — verification is the mechanism that makes the rest correct.

**Root cause.** `internal/merge/merge.go` is single-threaded on one revset:
```
  revset := changeName + "@"
  jj rebase -b <change>@ -d main      # follows only <change>@'s line
  jj new main <change>@ -m "merge"    # parent is <change>@ — empty if @ is empty
  jj bookmark set main -r @           # advance; sibling/parent work never seen
```
If `@` is empty, the merge is empty. If the work is on a sibling, it's invisible. And because merge rewrites `<change>@` from the main session, the spawned workspace is left stale either way — with no signal that the merge was good or bad.

## What Changes

The merge becomes a **verification-gated** pipeline. The workspace's lifecycle is gated on *proof the work landed*, which is the user's stance: **keep the workspace until we know for sure the merge was successful.**

- **Robust work definition (closes instances 2 & 3).** `merge` SHALL define "the workspace's work" as **all non-empty commits in `main..(all workspace heads)`** — including divergent siblings and `@-` — not just `<change>@`. The rebase/merge SHALL operate over that real frontier.
- **Pre-merge op snapshot (rse4 Tension-1).** Before rewriting anything, `merge` SHALL capture the current `jj op log` head as a recovery handle to surface on failure.
- **Changed-file capture (rse4 L2).** Before merge, `merge` SHALL record the set of files touched across the whole work frontier (union over siblings, not just `<change>@`).
- **Post-merge smoke test (rse4 L1+L2).** After merge, `merge` SHALL verify (L1) main actually gained changes when the workspace had work, and (L2) every captured file is present on main. L3 (content equivalence) is documented as a follow-up.
- **Verification-gated workspace lifecycle (closes instance 1).**
  - **Proven** (smoke test passes) → the workspace's job is done → `merge` forgets the workspace. No live pointer means **no staleness is possible**.
  - **Unproven** (smoke test fails) → **non-destructive**: keep the workspace intact (a stale pointer here is an acceptable recovery state), do **not** move on, emit a **loud structured warning** naming the expected/missing files and the recovery handle (`jj op restore <preMergeOp>`), and exit non-zero. No auto-rollback (rse4 pre-1.0 stance).

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `merge`: merge SHALL operate on the workspace's full work frontier (not only `<change>@`), SHALL verify the work landed on main via a post-merge smoke test, and SHALL gate the workspace's lifecycle (forget vs keep) on that proof — never leaving a spawned workspace silently stale on success, never reporting success on an empty or orphaned merge.

## Impact

- **Code**: `internal/merge/merge.go` (revset → work-frontier, op snapshot, file capture, smoke test, forget-on-success / keep-on-failure). May reuse `cleanup.forgetWorkspace`'s pattern. New scenarios in `internal/merge/merge_integration_test.go` plus two new test helpers.
- **Severity**: CRITICAL — silent loss of merged proposals + pervasive staleness on the no-AI core capability.
- **Relation**: extends ADR-010 (`advanceMainToHead`) and `rebase-before-merge`; does not revert them. Absorbs jjay-rse4 (smoke test).
- **ADRs**: ADR-013 (verification-gated merge: prove the work landed, gate workspace lifecycle on proof).
- **Beans**: jjay-q6ko → in-progress, linked here. jjay-rse4 → scrapped (folded into this proposal).
- **Bean task ref**: [jjay-q6ko](../../.beans/jjay-q6ko--merging-now-almost-always-give-a-stale-error.md)
